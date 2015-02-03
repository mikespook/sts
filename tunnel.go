package sts

import (
	"io/ioutil"
	"net"
	"sync"

	"github.com/mikespook/golib/log"
	"golang.org/x/crypto/ssh"
)

func New(config *Config) *Server {
	return &Server{config: config}
}

type Server struct {
	sync.RWMutex
	config   *Config
	listener net.Listener
}

func (srv *Server) sshConfig() (config *ssh.ServerConfig, err error) {
	config = &ssh.ServerConfig{
		NoClientAuth: false,
	}
	config.PasswordCallback = func(conn ssh.ConnMetadata, password []byte) (p *ssh.Permissions, err error) {
		return nil, nil
	}

	for _, key := range srv.config.Keys {
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

func (srv *Server) Serve() (err error) {
	srv.listener, err = net.Listen("tcp", srv.config.Addr)
	if err != nil {
		return
	}
	var sshConfig *ssh.ServerConfig
	sshConfig, err = srv.sshConfig()
	if err != nil {
		return
	}
	for {
		var conn net.Conn
		conn, err = srv.listener.Accept()
		if err != nil {
			log.Errorf("Accept: %s", err)
			continue
		}
		go doSession(conn, sshConfig)
	}

	return nil
}

func (srv *Server) Close() error {
	return srv.listener.Close()
}
