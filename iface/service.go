package iface

import "time"

type Service interface {
	Config(config interface{}) error

	Serve() error
	Close() error

	ETime() time.Time
}
