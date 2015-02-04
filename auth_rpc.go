package sts

import (
	"net/rpc"
	"net/url"

	"golang.org/x/crypto/ssh"
)

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
		client, err = rpc.Dial(u.Scheme, u.Host)
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
		return nil, ErrAuthFailed
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
		return nil, ErrAuthFailed
	}
}
