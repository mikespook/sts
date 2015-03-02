package main

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Addr string
	Log  struct {
		File, Level string
	}
	Mongo struct {
		Addr, Db string
	}
	Pwd string
}

func LoadConfig(filename string) (cfg *Config, err error) {
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
	return
}
