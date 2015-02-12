package sts

import (
	"net"
	"net/http"
	"net/rpc"
	"net/url"
)

func NewRPC(addr string) *RPC {
	return &RPC{
		addr:   addr,
		server: rpc.NewServer(),
	}
}

type RPC struct {
	addr string

	server   *rpc.Server
	listener net.Listener
}

func (srv *RPC) Serve() error {
	u, err := url.Parse(srv.addr)
	if err != nil {
		return err
	}
	isHttp := u.Scheme == "http"
	if isHttp {
		srv.server.HandleHTTP(u.Path, "/_debug")
	} else if u.Scheme == "" {
		u.Scheme = "tcp"
		u.Host = u.Path
		u.Path = ""
	}
	if srv.listener, err = net.Listen(u.Scheme, u.Host); err != nil {
		return err
	}
	if isHttp {
		http.Serve(srv.listener, nil)
	} else {
		srv.server.Accept(srv.listener)
	}
	return nil
}

func (srv *RPC) Close() error {
	return srv.listener.Close()
}
