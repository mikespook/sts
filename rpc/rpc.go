package rpc

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"net/url"

	"github.com/mikespook/sts/bus"
)

func New() bus.Service {
	return &RPC{
		Server: rpc.NewServer(),
	}
}

type RPC struct {
	*rpc.Server

	config   *Config
	listener net.Listener
	bus      bus.Bus
}

func (srv *RPC) Bus(bus bus.Bus) {
	srv.bus = bus
}

func (srv *RPC) Config(config interface{}) (err error) {
	cfg, ok := config.(*Config)
	if !ok {
		err = fmt.Errorf("Wrong paramater %T, wants %T", config, cfg)
		return
	}
	srv.config = cfg
	return
}

func (srv *RPC) Serve() error {
	if err := srv.RegisterName("Ctrl", &Ctrl{srv.bus.Ctrl()}); err != nil {
		return err
	}
	if err := srv.RegisterName("Stat", &Stat{srv.bus.Stat()}); err != nil {
		return err
	}
	u, err := url.Parse(srv.config.Addr)
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

func (srv *RPC) Restart() error {
	return nil
}
