package main

import (
	"net"
	"net/http"
	"net/rpc"
	"net/url"

	"github.com/mikespook/sts"
	"github.com/mikespook/sts/model"
	"golang.org/x/crypto/ssh"
	"gopkg.in/mgo.v2"
)

type STS struct {
	session *mgo.Session
	dbName  string
}

func (data *STS) PasswordAuth(auth *model.Auth, perm *ssh.Permissions) error {
	if auth.User == "" {
		return sts.ErrAuthFailed
	}
	if auth.Password == nil {
		return sts.ErrAuthFailed
	}
	user, err := model.GetUser(data.session, data.dbName, auth.User)
	if err != nil {
		return err
	}
	if user.CheckPassword(auth.Password) {
		*perm = user.Permissions
		return nil
	}
	return sts.ErrAuthFailed
}

func (data *STS) PublicKeyAuth(auth *model.Auth, perm *ssh.Permissions) error {
	if auth.User == "" {
		return sts.ErrAuthFailed
	}
	if auth.Key == nil {
		return sts.ErrAuthFailed
	}
	user, err := model.GetUser(data.session, data.dbName, auth.User)
	if err != nil {
		return err
	}
	if user.CheckPublicKey(auth.Key) {
		*perm = user.Permissions
		return nil
	}
	return sts.ErrAuthFailed
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
	if err := srv.server.RegisterName("STS", &STS{session, config.Mongo.Db}); err != nil {
		return nil, err
	}
	var u *url.URL
	if u, err = url.Parse(config.Addr); err != nil {
		return nil, err
	}
	if u.Scheme == "http" {
		srv.server.HandleHTTP(u.Path, "")
		srv.isHttp = true
	}
	if srv.listener, err = net.Listen("tcp", u.Host); err != nil {
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
