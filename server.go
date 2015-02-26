package sts

import (
	"fmt"
	"os"

	"github.com/mikespook/golib/log"
	"github.com/mikespook/sts/iface"
	"github.com/mikespook/sts/rpc"
	"github.com/mikespook/sts/tunnel"
)

const (
	Tunnel = "Tunnel"
	RPC    = "RPC"
)

func New(cfg *Config) *Server {
	srv := &Server{
		config:    cfg,
		errExit:   make(chan error),
		errCommon: make(chan error),
		services:  make(map[string]iface.Service),
	}
	return srv
}

type Server struct {
	services map[string]iface.Service

	errExit   chan error
	errCommon chan error

	config *Config
}

func (srv *Server) Serve() (err error) {
	log.Messagef("Set PWD: %s", srv.config.Pwd)
	if err = os.Chdir(srv.config.Pwd); err != nil {
		return
	}
	go srv.start(rpc.New, RPC, &srv.config.RPC)
	go srv.start(tunnel.New, Tunnel, &srv.config.Tunnel)
	return srv.wait()
}

func (srv *Server) Close() {
	srv.close(Tunnel)
	srv.close(RPC)
	srv.shutdown()
}

func (srv *Server) reboot() {
	srv.close(Tunnel)
	go srv.start(tunnel.New, Tunnel, srv.config.Tunnel)
}

func (srv *Server) wait() (err error) {
Loop:
	for {
		select {
		case err = <-srv.errExit:
			break Loop
		case err = <-srv.errCommon:
			log.Error(err)
		}
	}
	return
}

func (srv *Server) shutdown() {
	close(srv.errExit)
	close(srv.errCommon)
}

func (srv *Server) start(f func() iface.Service, name string, config interface{}) {
	log.Messagef("Start %s: %+v", name, config)
	service := f()
	if err := service.Config(config); err != nil {
		srv.errExit <- fmt.Errorf("%s Start: %s", name, err)
		return
	}
	if err := service.Serve(); err != nil {
		srv.errExit <- fmt.Errorf("%s Serve: %s", name, err)
		return
	}
	srv.services[name] = service
}

func (srv *Server) close(name string) {
	log.Messagef("Close %s", name)
	service, ok := srv.services[name]
	if !ok {
		return
	}
	if err := service.Close(); err != nil {
		srv.errCommon <- fmt.Errorf("%s Close: %s", name, err)
	}
}

func (srv *Server) Ctrl() iface.Ctrl {
	return nil
}

func (srv *Server) Stat() iface.Stat {
	return nil
}
