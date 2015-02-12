package sts

import (
	"io/ioutil"
	"os"

	"github.com/mikespook/sts/auth"
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
	auth *auth.Config
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
	cfg.auth, err = auth.LoadConfig(cfg.Auth)
	return
}
