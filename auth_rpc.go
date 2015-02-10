package sts

import (
	"net/rpc"
	"net/url"
	"github.com/mikespook/sts/model"
	"golang.org/x/crypto/ssh"
)

const (
	RPCPasswordAuth = "STS.PasswordAuth"
	RPCPubKeyAuth   = "STS.PublicKeyAuth"
)

func init() {
	RegisterAuth(AuthPassword, AuthRPC, rpcPasswordHandle)
	RegisterAuth(AuthPubKey, AuthRPC, rpcPubKeyHandle)
}

func rpcPasswordHandle(cfg *configAuth, key, prefix, value string) (exclusive bool, err error) {
	if cfg.Password, err = newRpcPasswordAuth(value); err != nil {
		return false, err
	}
	return false, nil
}

func rpcPubKeyHandle(cfg *configAuth, key, prefix, value string) (exclusive bool, err error) {
	if cfg.PublicKey, err = newRpcPublicKeyAuth(value); err != nil {
		return false, err
	}
	return false, nil
}

func newRpcClient(rawurl string) (client *rpc.Client, err error) {
	var u *url.URL
	if u, err = url.Parse(rawurl); err != nil {
		return
	}
	if u.Scheme == "http" {
		if u.Path == "" {
			client, err = rpc.DialHTTP("tcp", u.Host)
		} else {
			client, err = rpc.DialHTTPPath("tcp", u.Host, u.Path)
		}
	} else {
		if u.Scheme == "" {
			u.Scheme = "tcp"
		}
		client, err = rpc.Dial(u.Scheme, u.Path)
	}
	return
}

func newRpcPasswordAuth(rawurl string) (PasswordAuth, error) {
	client, err := newRpcClient(rawurl)
	if err != nil {
		return nil, err
	}
	return &rpcPasswordAuth{client}, nil
}

type rpcPasswordAuth struct {
	client *rpc.Client
}

func (a *rpcPasswordAuth) Callback() passwordCallback {
	return func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
		args := &model.Auth{
			Addr:     conn.RemoteAddr().String(),
			User:     conn.User(),
			Password: model.Password(password),
		}
		perm := &ssh.Permissions{}
		if err := a.client.Call(RPCPasswordAuth, args, &perm); err != nil {
			return nil, err
		}
		return perm, nil
	}
}

func newRpcPublicKeyAuth(rawurl string) (PublicKeyAuth, error) {
	client, err := newRpcClient(rawurl)
	if err != nil {
		return nil, err
	}
	return &rpcPublicKeyAuth{client}, nil
}

type rpcPublicKeyAuth struct {
	client *rpc.Client
}

func (a *rpcPublicKeyAuth) Callback() publicKeyCallback {
	return func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
		args := &model.Auth{
			Addr: conn.RemoteAddr().String(),
			User: conn.User(),
			Key:  key.Marshal(),
		}
		perm := &ssh.Permissions{}
		if err := a.client.Call(RPCPubKeyAuth, args, &perm); err != nil {
			return nil, err
		}
		return perm, nil
	}
}
