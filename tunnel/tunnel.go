package tunnel

import (
	"fmt"
	"io/ioutil"
	"net"

	"github.com/mikespook/golib/log"
	"github.com/mikespook/sts/model"
	"github.com/mikespook/sts/tunnel/auth"
	"golang.org/x/crypto/ssh"
)

func New(states model.States) model.Service {
	return &Tunnel{
		states: states,
	}
}

type Tunnel struct {
	config   *Config
	auth     *auth.Config
	listener net.Listener
	states   model.States
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

	for _, key := range tun.config.Keys {
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

func (tun *Tunnel) Config(config interface{}) (err error) {
	cfg, ok := config.(*Config)
	if !ok {
		err = fmt.Errorf("Wrong paramater %t, wants %t", config, cfg)
		return
	}
	tun.config = cfg
	tun.auth, err = auth.LoadConfig(cfg.Auth)
	return
}

func (tun *Tunnel) Serve() (err error) {
	tun.listener, err = net.Listen("tcp", tun.config.Addr)
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
			if opErr, ok := err.(*net.OpError); ok {
				if opErr.Temporary() || opErr.Timeout() {
					log.Errorf("Accept: %s", opErr)
					continue
				}
				break
			}
			log.Errorf("Accept: %s", err)
			break
		}
		go tun.session(conn, sshConfig)
	}

	return nil
}

func (tun *Tunnel) Close() error {
	//	tun.state.Close()
	return tun.listener.Close()
}

func (tun *Tunnel) Restart() error {
	return nil
}

func (tun *Tunnel) session(conn net.Conn, config *ssh.ServerConfig) {
	defer func() {
		conn.Close()
		log.Messagef("Disconnect: %s", conn.RemoteAddr())
	}()
	log.Messagef("Connect: %s", conn.RemoteAddr())
	s, err := newSession(conn, config, tun.states)
	if err != nil {
		log.Errorf("SSH-Connect: %s", err)
		return
	}
	tun.states.Sessions().Add(s)
	defer tun.states.Sessions().Remove(s)
	defer s.Close()
	s.Serve()
}
