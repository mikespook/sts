package sts

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type config struct {
	Pwd  string
	Addr string
	Keys []string
	Log  struct {
		File, Level string
	}
	Auth map[string]string
	auth configAuth
}

type configAuth struct {
	anonymous bool
	PublicKey PublicKeyAuth
	Password  PasswordAuth
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
	return parseAuth(&cfg.auth, cfg.Auth)
}
