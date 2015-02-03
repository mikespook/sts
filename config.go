package sts

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Addr string
	Keys []string
	Log  struct {
		File, Level string
	}
}

func LoadConfig(filename string) (config *Config, err error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(data, &config)
	return
}
