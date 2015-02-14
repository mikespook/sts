package sts

import (
	"fmt"
	"net"
	"os"

	"github.com/mikespook/golib/log"
	"github.com/mikespook/sts/bus"
	"github.com/mikespook/sts/rpc"
	"github.com/mikespook/sts/tunnel"
	"golang.org/x/crypto/ssh"
	"gopkg.in/mgo.v2/bson"
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
		services:  make(map[string]bus.Service),
		sessions:  make(map[bson.ObjectId]*session),
	}
	return srv
}

type Server struct {
	services map[string]bus.Service

	errExit   chan error
	errCommon chan error

	config *Config

	sessions map[bson.ObjectId]*session
}

func (srv *Server) Session(conn net.Conn, config *ssh.ServerConfig) {
	// TODO register session
	defer func() {
		conn.Close()
		log.Messagef("Disconnect: %s", conn.RemoteAddr())
	}()
	log.Messagef("Connect: %s", conn.RemoteAddr())

	s, err := newSession(conn, config)
	if err != nil {
		log.Errorf("SSH-Connect: %s", err)
		return
	}
	log.Messagef("SSH-Connect: %s [%s@%s] (%s)", s.Id.Hex(), s.SshConn.User(),
		s.SshConn.RemoteAddr(), s.SshConn.ClientVersion())
	// TODO need a lock or use another data structure to mantain the list
	srv.sessions[s.Id] = s
	s.Do()
	log.Messagef("SSH-Disconnect: %s", s.Id.Hex())
}

func (srv *Server) Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func (srv *Server) Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func (srv *Server) Warningf(format string, args ...interface{}) {
	log.Warningf(format, args...)
}

func (srv *Server) Messagef(format string, args ...interface{}) {
	log.Messagef(format, args...)
}

func (srv *Server) Serve() (err error) {
	srv.Messagef("Set PWD: %s", srv.config.Pwd)
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

func (srv *Server) start(f func() bus.Service, name string, config interface{}) {
	srv.Messagef("Start %s: %+v", name, config)
	service := f()
	service.Bus(srv)
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
	srv.Messagef("Close %s", name)
	service, ok := srv.services[name]
	if !ok {
		return
	}
	if err := service.Close(); err != nil {
		srv.errCommon <- fmt.Errorf("%s Close: %s", name, err)
	}
}

func (srv *Server) Ctrl() bus.Ctrl {
	return nil
}

func (srv *Server) Stat() bus.Stat {
	return nil
}
