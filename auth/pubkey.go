package auth

import (
	"bytes"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

func init() {
	Register(KeyPubKey, PrefixStatic, staticPubKeyHandle)
	Register(KeyPubKey, PrefixFile, filePubKeyHandle)
}

func staticPubKeyHandle(cfg *Config, key, prefix, value string) (exclusive bool, err error) {
	if cfg.PublicKey, err = newStaticPublicKey([]byte(value)); err != nil {
		return false, err
	}
	return false, nil
}

func filePubKeyHandle(cfg *Config, key, prefix, value string) (exclusive bool, err error) {
	if cfg.PublicKey, err = newFilePublicKey(value); err != nil {
		return false, err
	}
	return false, nil
}

type staticPublicKey struct {
	keyBytes []byte
}

// The public key read from the config field
func newStaticPublicKey(key []byte) (PublicKey, error) {
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
		return nil, ErrFailed
	}
}

// The public key read from the file `f`
func newFilePublicKey(f string) (PublicKey, error) {
	key, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return newStaticPublicKey(key)
}
