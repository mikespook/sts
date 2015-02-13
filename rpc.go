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
		Server: rpc.NewServer(),
	}
}

type RPC struct {
	*rpc.Server

	addr     string
	listener net.Listener
	bus      Bus
}

func (srv *RPC) Serve() error {
	if err := srv.RegisterName("Ctrl", &Ctrl{srv.bus.Server()}); err != nil {
		return err
	}
	if err := srv.RegisterName("Stat", &Stat{srv.bus.Server()}); err != nil {
		return err
	}
	u, err := url.Parse(srv.addr)
	if err != nil {
		return err
	}
	isHttp := u.Scheme == "http"
	if isHttp {
		srv.HandleHTTP(u.Path, "/_debug")
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
		srv.Accept(srv.listener)
	}
	return nil
}

func (srv *RPC) Close() error {
	return srv.listener.Close()
}
