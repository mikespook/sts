package sts

import (
	"bytes"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

type passwordCallback func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error)

// Password callback
type PasswordAuth interface {
	Callback() passwordCallback
}

// The password read from the config field
func newStaticPassword(password []byte) PasswordAuth {
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
		return nil, ErrAuthFailed
	}
}

// The password read from the file `f`
func newFilePassword(f string) (PasswordAuth, error) {
	password, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return newStaticPassword(password), nil
}
