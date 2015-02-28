package tunnel

import (
	"fmt"
	"io/ioutil"
	"net"
	"time"

	"github.com/mikespook/golib/log"
	"github.com/mikespook/sts/iface"
	"github.com/mikespook/sts/tunnel/auth"
	"golang.org/x/crypto/ssh"
)

func New(keeper iface.Keeper) iface.Service {
	return &Tunnel{
		keeper: keeper,
	}
}

type Tunnel struct {
	config   *Config
	auth     *auth.Config
	listener net.Listener
	keeper   iface.Keeper
	etime    time.Time
}

func (tun *Tunnel) ETime() time.Time {
	return tun.etime
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
	tun.etime = time.Now()
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
	return tun.listener.Close()
}

func (tun *Tunnel) session(conn net.Conn, config *ssh.ServerConfig) {
	defer func() {
		conn.Close()
		log.Messagef("Disconnect: %s", conn.RemoteAddr())
	}()
	log.Messagef("Connect: %s", conn.RemoteAddr())
	s, err := newSession(conn, config)
	if err != nil {
		log.Errorf("SSH-Connect: %s", err)
		return
	}
	s.keeper = tun.keeper
	tun.keeper.AddSession(s)
	defer tun.keeper.RemoveSession(s)
	s.Serve()
}
