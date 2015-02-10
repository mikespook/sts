package sts

import (
	"errors"
	"strings"
)

const (
	AuthStatic      = "static"
	AuthFile        = "file"
	AuthRPC         = "rpc"
	AuthAnonymous   = "anonymous"
	AuthPassword    = "password"
	AuthPubKey      = "pubkey"
	UnderstandRasks = "I Understand the Risks"
)

type authHandler func(cfg *configAuth, key, prefix, value string) (bool, error)

var (
	// Callbacks should retuen this error, whatever the real reason is.
	ErrAuthFailed = errors.New("Auth failed")
	auths         = make(map[string]authHandler)
)

func RegisterAuth(key, prefix string, f authHandler) {
	auths[key+"&"+prefix] = f
}

func parseAuth(cfg *configAuth, data map[string]string) error {
	for key, value := range data {
		tmp := strings.SplitN(value, "://", 2)
		if len(tmp) == 2 {
			if exclusive, err := auths[key+"&"+tmp[0]](cfg, key, tmp[0], tmp[1]); err != nil {
				return err
			} else if exclusive {
				break
			}
		}
	}
	return nil
}

func init() {
	RegisterAuth(AuthAnonymous, "", anonymousHandle)
}

func anonymousHandle(cfg *configAuth, key, prefix, value string) (bool, error) {
	if value == UnderstandRasks {
		cfg.anonymous = true
	}
	return true, nil
}
