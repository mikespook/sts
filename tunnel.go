package sts

import (
	"io/ioutil"
	"net"

	"github.com/mikespook/golib/log"
	"github.com/mikespook/sts/auth"
	"golang.org/x/crypto/ssh"
)

type Tunnel struct {
	// Config
	addr string
	keys []string
	auth *auth.Config

	listener net.Listener
}

func NewTunnel(addr string, keys []string, config *auth.Config) *Tunnel {
	return &Tunnel{
		addr: addr,
		keys: keys,
		auth: config,
	}
}

func (tun *Tunnel) sshConfig() (config *ssh.ServerConfig, err error) {
	config = &ssh.ServerConfig{
		NoClientAuth: tun.auth.Anonymous,
	}
	if !tun.auth.Anonymous {
		if tun.auth.Password != nil {
			config.PasswordCallback = tun.auth.Password.Callback()
		}
		if tun.auth.PublicKey != nil {
			config.PublicKeyCallback = tun.auth.PublicKey.Callback()
		}
	}

	for _, key := range tun.keys {
		var privBytes []byte
		if privBytes, err = ioutil.ReadFile(key); err != nil {
			return
		}
		var privKey ssh.Signer
		privKey, err = ssh.ParsePrivateKey(privBytes)
		if err != nil {
			return
		}
		config.AddHostKey(privKey)
	}
	return
}

func (tun *Tunnel) Serve() (err error) {
	tun.listener, err = net.Listen("tcp", tun.addr)
	if err != nil {
		return
	}
	var sshConfig *ssh.ServerConfig
	sshConfig, err = tun.sshConfig()
	if err != nil {
		return
	}
	for {
		var conn net.Conn
		conn, err = tun.listener.Accept()
		if err != nil {
			log.Errorf("Accept: %s", err)
			continue
		}
		go doSession(conn, sshConfig)
	}

	return nil
}

func (tun *Tunnel) Close() error {
	return tun.listener.Close()
}
