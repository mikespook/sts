package auth

import (
	"bytes"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

func init() {
	Register(KeyPassword, PrefixStatic, staticPasswordHandle)
	Register(KeyPassword, PrefixFile, filePasswordHandle)
}

func staticPasswordHandle(cfg *Config, key, prefix, value string) (bool, error) {
	cfg.Password = newStaticPassword([]byte(value))
	return false, nil
}

func filePasswordHandle(cfg *Config, key, prefix, value string) (exclusive bool, err error) {
	if cfg.Password, err = newFilePassword(value); err != nil {
		return false, err
	}
	return false, nil
}

// The password read from the config field
func newStaticPassword(password []byte) Password {
	return &staticPassword{password}
}

type staticPassword struct {
	password []byte
}

func (sp *staticPassword) Callback() passwordCallback {
	return func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
		if bytes.Compare(password, sp.password) == 0 {
			return nil, nil
		}
		return nil, ErrFailed
	}
}

// The password read from the file `f`
func newFilePassword(f string) (Password, error) {
	password, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return newStaticPassword(password), nil
}
