package auth

import (
	"net/rpc"
	"net/url"

	"github.com/mikespook/sts/model"
	"golang.org/x/crypto/ssh"
)

const (
	RPCPassword = "Auth.Password"
	RPCPubKey   = "Auth.PublicKey"
)

func init() {
	Register(KeyPassword, PrefixRPC, rpcPasswordHandle)
	Register(KeyPubKey, PrefixRPC, rpcPubKeyHandle)
}

func rpcPasswordHandle(cfg *Config, key, prefix, value string) (exclusive bool, err error) {
	if cfg.Password, err = newRpcPassword(value); err != nil {
		return false, err
	}
	return false, nil
}

func rpcPubKeyHandle(cfg *Config, key, prefix, value string) (exclusive bool, err error) {
	if cfg.PublicKey, err = newRpcPublicKey(value); err != nil {
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

func newRpcPassword(rawurl string) (Password, error) {
	client, err := newRpcClient(rawurl)
	if err != nil {
		return nil, err
	}
	return &rpcPassword{client}, nil
}

type rpcPassword struct {
	client *rpc.Client
}

func (a *rpcPassword) Callback() passwordCallback {
	return func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
		args := &model.Auth{
			Addr:     conn.RemoteAddr().String(),
			User:     conn.User(),
			Password: model.HashPassword(password),
		}
		perm := &ssh.Permissions{}
		if err := a.client.Call(RPCPassword, args, &perm); err != nil {
			return nil, err
		}
		return perm, nil
	}
}

func newRpcPublicKey(rawurl string) (PublicKey, error) {
	client, err := newRpcClient(rawurl)
	if err != nil {
		return nil, err
	}
	return &rpcPublicKey{client}, nil
}

type rpcPublicKey struct {
	client *rpc.Client
}

func (a *rpcPublicKey) Callback() publicKeyCallback {
	return func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
		args := &model.Auth{
			Addr: conn.RemoteAddr().String(),
			User: conn.User(),
			Key:  key.Marshal(),
		}
		perm := &ssh.Permissions{}
		if err := a.client.Call(RPCPubKey, args, &perm); err != nil {
			return nil, err
		}
		return perm, nil
	}
}
