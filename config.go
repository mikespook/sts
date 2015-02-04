package sts

import (
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	AuthStatic      = "static://"
	AuthFile        = "file://"
	AuthRPC         = "rpc://"
	AuthAnonymous   = "anonymous"
	AuthPassword    = "password"
	AuthPubKey      = "pubkey"
	UnderstandRasks = "I Understand the Risks"
)

type config struct {
	Pwd  string
	Addr string
	Keys []string
	Log  struct {
		File, Level string
	}
	Auth map[string]string
	auth struct {
		anonymous bool
		PublicKey PublicKeyAuth
		Password  PasswordAuth
	}
}

func LoadConfig(filename string) (cfg *config, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return
	}
	if err = os.Chdir(cfg.Pwd); err != nil {
		return
	}
	err = cfg.parseAuth()
	return
}

func (cfg *config) parseAuth() (err error) {
Loop:
	for key, auth := range cfg.Auth {
		switch key {
		case AuthAnonymous:
			switch strings.ToLower(auth) {
			case UnderstandRasks:
				cfg.auth.anonymous = true
				break Loop
			}
		case AuthPassword:
			switch {
			case strings.HasPrefix(auth, AuthStatic):
				cfg.auth.Password = newStaticPassword([]byte(strings.TrimPrefix(auth, AuthStatic)))
			case strings.HasPrefix(auth, AuthFile):
				if cfg.auth.Password, err = newFilePassword(strings.TrimPrefix(auth, AuthFile)); err != nil {
					return
				}
			case strings.HasPrefix(auth, AuthRPC):
				if cfg.auth.Password, err = newRpcPasswordAuth(strings.TrimPrefix(auth, AuthRPC)); err != nil {
					return
				}
			}
		case AuthPubKey:
			switch {
			case strings.HasPrefix(auth, AuthStatic):
				if cfg.auth.PublicKey, err = newStaticPublicKey([]byte(strings.TrimPrefix(auth, AuthStatic))); err != nil {
					return
				}
			case strings.HasPrefix(auth, AuthFile):
				if cfg.auth.PublicKey, err = newFilePublicKey(strings.TrimPrefix(auth, AuthFile)); err != nil {
					return
				}
			case strings.HasPrefix(auth, AuthRPC):
				if cfg.auth.PublicKey, err = newRpcPublicKeyAuth(strings.TrimPrefix(auth, AuthRPC)); err != nil {
					return
				}
			}
		}
	}
	return
}
