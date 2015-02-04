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

type Config struct {
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

func LoadConfig(filename string) (config *Config, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	if err = yaml.Unmarshal(data, &config); err != nil {
		return
	}
	if err = os.Chdir(config.Pwd); err != nil {
		return
	}
	err = config.parseAuth()
	return
}

func (config *Config) parseAuth() (err error) {
Loop:
	for key, auth := range config.Auth {
		switch key {
		case AuthAnonymous:
			switch strings.ToLower(auth) {
			case UnderstandRasks:
				config.auth.anonymous = true
				break Loop
			}
		case AuthPassword:
			switch {
			case strings.HasPrefix(auth, AuthStatic):
				config.auth.Password = NewStaticPassword([]byte(strings.TrimPrefix(auth, AuthStatic)))
			case strings.HasPrefix(auth, AuthFile):
				if config.auth.Password, err = NewFilePassword(strings.TrimPrefix(auth, AuthFile)); err != nil {
					return
				}
			case strings.HasPrefix(auth, AuthRPC):
				// TODO
			}
		case AuthPubKey:
			switch {
			case strings.HasPrefix(auth, AuthStatic):
				if config.auth.PublicKey, err = NewStaticPublicKey([]byte(strings.TrimPrefix(auth, AuthStatic))); err != nil {
					return
				}
			case strings.HasPrefix(auth, AuthFile):
				if config.auth.PublicKey, err = NewFilePublicKey(strings.TrimPrefix(auth, AuthFile)); err != nil {
					return
				}
			case strings.HasPrefix(auth, AuthRPC):
				// TODO
			}
		}
	}
	return
}
