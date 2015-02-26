package rpc

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"net/url"

	"github.com/mikespook/sts/iface"
)

func New() iface.Service {
	return &RPC{
		Server: rpc.NewServer(),
	}
}

type RPC struct {
	*rpc.Server

	config   *Config
	listener net.Listener
	daemon   iface.Daemon
}

func (srv *RPC) Daemon(daemon iface.Daemon) {
	srv.daemon = daemon
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
	if err := srv.RegisterName("Ctrl", &Ctrl{srv.daemon.Ctrl()}); err != nil {
		return err
	}
	if err := srv.RegisterName("Stat", &Stat{srv.daemon.Stat()}); err != nil {
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
