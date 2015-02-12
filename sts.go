package sts

import (
	"sync"

	"github.com/mikespook/golib/log"
)

func New(cfg *config) *Server {
	return &Server{
		config: cfg,
		rpc:    NewRPC(cfg.Addr.RPC),
		tunnel: NewTunnel(cfg.Addr.Tunnel, cfg.Keys, cfg.auth),
	}
}

type Server struct {
	sync.WaitGroup

	config *config
	status *Status

	rpc    *RPC
	tunnel *Tunnel
}

func (srv *Server) Serve() {
	go func() {
		srv.Add(1)
		defer srv.Done()
		if err := srv.tunnel.Serve(); err != nil {
			log.Errorf("Tunnel Serve: %s", err)
		}
	}()
	go func() {
		srv.Add(1)
		defer srv.Done()
		if err := srv.rpc.Serve(); err != nil {
			log.Errorf("RPC Serve: %s", err)
		}
	}()
	srv.Wait()
}

func (srv *Server) Close() {
	if err := srv.tunnel.Close(); err != nil {
		log.Errorf("Tunnel Close: %s", err)
	}
	if err := srv.rpc.Close(); err != nil {
		log.Errorf("RPC Close: %s", err)
	}
	return
}
