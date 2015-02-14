package tunnel

import (
	"fmt"
	"io/ioutil"
	"net"

	"github.com/mikespook/sts/auth"
	"github.com/mikespook/sts/bus"
	"golang.org/x/crypto/ssh"
)

func New() bus.Service {
	return &Tunnel{}
}

type Tunnel struct {
	config   *Config
	auth     *auth.Config
	listener net.Listener

	bus bus.Bus
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

func (tun *Tunnel) Bus(bus bus.Bus) {
	tun.bus = bus
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
					tun.bus.Errorf("Accept: %s", opErr)
					continue
				}
				break
			}
			tun.bus.Errorf("Accept: %s", err)
			break
		}
		go tun.bus.Session(conn, sshConfig)
	}

	return nil
}

func (tun *Tunnel) Close() error {
	tun.bus.Close()
	return tun.listener.Close()
}

func (tun *Tunnel) Restart() error {
	return nil
}
