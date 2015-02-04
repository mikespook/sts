package sts

import (
	"bytes"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

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
