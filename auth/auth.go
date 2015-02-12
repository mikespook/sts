package auth

import (
	"errors"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
)

const (
	PrefixStatic = "static"
	PrefixFile   = "file"
	PrefixRPC    = "rpc"

	KeyAnonymous = "anonymous"
	KeyPassword  = "password"
	KeyPubKey    = "pubkey"

	ValueUnderstandRasks = "I Understand the Risks"
)

type passwordCallback func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error)

// Password callback
type Password interface {
	Callback() passwordCallback
}

type publicKeyCallback func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error)

// PublicKey callback
type PublicKey interface {
	Callback() publicKeyCallback
}

// Configuration of Authorization
type Config struct {
	Anonymous bool
	PublicKey PublicKey
	Password  Password
}

// Return `ture` for exclusive options
type OptionsHandler func(cfg *Config, key, prefix, value string) (bool, error)

var (
	// Callbacks should retuen this error, whatever the real reason is.
	ErrFailed = errors.New("Auth failed")
	options   = make(map[string]OptionsHandler)
	mutex     sync.Mutex
)

func Register(key, prefix string, f OptionsHandler) {
	mutex.Lock()
	options[key+"-"+prefix] = f
	mutex.Unlock()
}

func LoadConfig(data map[string]string) (cfg *Config, err error) {
	cfg = &Config{}
	for key, value := range data {
		tmp := strings.SplitN(value, "://", 2)
		if len(tmp) == 2 {
			if exclusive, err := options[key+"-"+tmp[0]](cfg, key, tmp[0], tmp[1]); err != nil {
				return nil, err
			} else if exclusive {
				break
			}
		}
	}
	return cfg, nil
}

func init() {
	Register(KeyAnonymous, "", anonymousHandle)
}

func anonymousHandle(cfg *Config, key, prefix, value string) (bool, error) {
	if value == ValueUnderstandRasks {
		cfg.Anonymous = true
	}
	return true, nil
}
