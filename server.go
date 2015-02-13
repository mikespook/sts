// TODO refactoring Bus design
package sts

import (
	"fmt"
	"net"

	"github.com/mikespook/golib/log"
	"github.com/mikespook/sts/auth"
	"golang.org/x/crypto/ssh"
	"gopkg.in/mgo.v2/bson"
)

type Bus interface {
	Server() *Server
	Errorf(format string, args ...interface{})
	Messagef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Warningf(format string, args ...interface{})

	Session(conn net.Conn, config *ssh.ServerConfig)
}

func New(cfg *Config) *Server {
	srv := &Server{
		config:    cfg,
		errExit:   make(chan error),
		errCommon: make(chan error),
		sessions:  make(map[bson.ObjectId]*session),
	}
	return srv
}

type Server struct {
	rpc    *RPC
	tunnel *Tunnel

	errExit   chan error
	errCommon chan error

	config *Config
	status *Status

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

func (srv *Server) Server() *Server {
	return srv
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
	go srv.startTunnel(srv.config.Addr.Tunnel, srv.config.Keys, srv.config.auth)
	go srv.startRPC(srv.config.Addr.RPC)
	return srv.wait()
}

func (srv *Server) Close() error {
	srv.closeTunnel()
	srv.closeRPC()
	srv.shutdown()
	return nil
}

func (srv *Server) reboot() {
	srv.closeTunnel()
	go srv.startTunnel(srv.config.Addr.Tunnel, srv.config.Keys, srv.config.auth)
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

func (srv *Server) startTunnel(addr string, keys []string, config *auth.Config) {
	srv.tunnel = NewTunnel(addr, keys, config)
	srv.tunnel.bus = srv
	if err := srv.tunnel.Serve(); err != nil {
		srv.errCommon <- fmt.Errorf("Tunnel Serve: %s", err)
	}
}

func (srv *Server) closeTunnel() {
	if err := srv.tunnel.Close(); err != nil {
		srv.errCommon <- fmt.Errorf("Tunnel Close: %s", err)
	}
}

func (srv *Server) startRPC(addr string) {
	srv.rpc = NewRPC(addr)
	srv.rpc.bus = srv
	if err := srv.rpc.Serve(); err != nil {
		srv.errExit <- fmt.Errorf("RPC Serve: %s", err)
	}
}

func (srv *Server) closeRPC() {
	if err := srv.rpc.Close(); err != nil {
		srv.errCommon <- fmt.Errorf("RPC Close: %s", err)
	}
}

func (srv *Server) kickoff(id bson.ObjectId) error {
	return nil
}

func (srv *Server) cutoff(id bson.ObjectId) error {
	return nil
}
