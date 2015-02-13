package sts

import (
	"io/ioutil"
	"os"

	"github.com/mikespook/sts/auth"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Pwd  string
	Addr struct {
		Tunnel string
		RPC    string
	}
	Keys []string
	Log  struct {
		File, Level string
	}
	Auth map[string]string
	auth *auth.Config
}

func LoadConfig(filename string) (cfg *Config, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return
	}
	err = cfg.Init()
	return
}

func (cfg *Config) Init() (err error) {
	if err = os.Chdir(cfg.Pwd); err != nil {
		return
	}
	cfg.auth, err = auth.LoadConfig(cfg.Auth)
	return
}
