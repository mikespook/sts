package main

import (
	"net"
	"net/http"
	"net/rpc"
	"net/url"

	"github.com/mikespook/sts/auth"
	"github.com/mikespook/sts/model"
	"golang.org/x/crypto/ssh"
	"gopkg.in/mgo.v2"
)

type Auth struct {
	session *mgo.Session
	dbName  string
}

func (a *Auth) Password(data *model.Auth, perm *ssh.Permissions) error {
	if data.User == "" {
		return auth.ErrFailed
	}
	if data.Password == nil {
		return auth.ErrFailed
	}
	user, err := model.GetUser(a.session, a.dbName, data.User)
	if err != nil {
		return err
	}
	if user.CheckPassword(data.Password) {
		*perm = user.Permissions
		return nil
	}
	return auth.ErrFailed
}

func (a *Auth) PublicKey(data *model.Auth, perm *ssh.Permissions) error {
	if data.User == "" {
		return auth.ErrFailed
	}
	if data.Key == nil {
		return auth.ErrFailed
	}
	user, err := model.GetUser(a.session, a.dbName, data.User)
	if err != nil {
		return err
	}
	if user.CheckPublicKey(data.Key) {
		*perm = user.Permissions
		return nil
	}
	return auth.ErrFailed
}

func NewRPC(config *Config) (*RPC, error) {
	session, err := mgo.Dial(config.Mongo.Addr)
	if err != nil {
		return nil, err
	}
	srv := &RPC{
		session: session,
		server:  rpc.NewServer(),
	}
	if err := srv.server.Register(&Auth{session, config.Mongo.Db}); err != nil {
		return nil, err
	}
	var u *url.URL
	if u, err = url.Parse(config.Addr); err != nil {
		return nil, err
	}
	if u.Scheme == "http" {
		srv.server.HandleHTTP(u.Path, "/_debug")
		srv.isHttp = true
	} else if u.Scheme == "" {
		u.Scheme = "tcp"
		u.Host = u.Path
		u.Path = ""
	}
	if srv.listener, err = net.Listen(u.Scheme, u.Host); err != nil {
		return nil, err
	}
	return srv, nil
}

type RPC struct {
	session  *mgo.Session
	server   *rpc.Server
	listener net.Listener
	isHttp   bool
}

func (srv *RPC) Serve() error {
	if srv.isHttp {
		http.Serve(srv.listener, nil)
	} else {
		srv.server.Accept(srv.listener)
	}
	return nil
}

func (srv *RPC) Close() error {
	srv.session.Close()
	return srv.listener.Close()
}
