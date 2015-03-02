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

func New(cfg *Config) *Sts {
	srv := &Sts{
		config:    cfg,
		errExit:   make(chan error),
		errCommon: make(chan error),
		services:  make(map[string]iface.Service),
	}
	srv.bus = newBus(srv)
	return srv
}

type Sts struct {
	services  map[string]iface.Service
	errExit   chan error
	errCommon chan error

	config *Config

	bus iface.Bus
}

func (srv *Sts) Serve() (err error) {
	log.Messagef("Set PWD: %s", srv.config.Pwd)
	if err = os.Chdir(srv.config.Pwd); err != nil {
		return
	}
	go srv.start(rpc.New, RPC, &srv.config.RPC)
	go srv.start(tunnel.New, Tunnel, &srv.config.Tunnel)
	return srv.wait()
}

func (srv *Sts) Close() {
	srv.close(Tunnel)
	srv.close(RPC)
	srv.shutdown()
}

func (srv *Sts) restart() {
	srv.close(Tunnel)
	go srv.start(tunnel.New, Tunnel, srv.config.Tunnel)
}

func (srv *Sts) wait() (err error) {
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

func (srv *Sts) shutdown() {
	close(srv.errExit)
	close(srv.errCommon)
}

func (srv *Sts) start(f func(iface.Bus) iface.Service, name string, config interface{}) {
	log.Messagef("Start %s: %+v", name, config)
	service := f(srv.bus)
	if err := service.Config(config); err != nil {
		srv.errExit <- fmt.Errorf("%s Start: %s", name, err)
		return
	}
	srv.services[name] = service
	if err := service.Serve(); err != nil {
		srv.errExit <- fmt.Errorf("%s Serve: %s", name, err)
		return
	}
}

func (srv *Sts) close(name string) {
	log.Messagef("Close %s", name)
	service, ok := srv.services[name]
	if !ok {
		return
	}
	if err := service.Close(); err != nil {
		srv.errCommon <- fmt.Errorf("%s Close: %s", name, err)
	}
}
