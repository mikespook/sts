package sts

import (
	"io/ioutil"

	"github.com/mikespook/sts/rpc"
	"github.com/mikespook/sts/tunnel"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Log struct {
		File, Level string
	}
	Pwd    string
	RPC    rpc.Config
	Tunnel tunnel.Config
}

func LoadConfig(filename string) (cfg *Config, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	err = yaml.Unmarshal(data, &cfg)
	return
}
