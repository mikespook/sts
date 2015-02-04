package sts

import (
	"bytes"
	"errors"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

// Callbacks should retuen this error, whatever the real reason is.
var ErrAuthFailed = errors.New("Auth failed")

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

type publicKeyCallback func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error)

// PublicKey callback
type PublicKeyAuth interface {
	Callback() publicKeyCallback
}

type staticPublicKey struct {
	keyBytes []byte
}

// The public key read from the config field
func newStaticPublicKey(key []byte) (PublicKeyAuth, error) {
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(key)
	if err != nil {
		return nil, err
	}
	return &staticPublicKey{pubKey.Marshal()}, nil
}

func (sp *staticPublicKey) Callback() publicKeyCallback {
	return func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
		if bytes.Compare(sp.keyBytes, key.Marshal()) == 0 {
			return nil, nil
		}
		return nil, ErrAuthFailed
	}
}

// The public key read from the file `f`
func newFilePublicKey(f string) (PublicKeyAuth, error) {
	key, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return newStaticPublicKey(key)
}
