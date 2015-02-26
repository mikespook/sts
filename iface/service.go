package iface

type Service interface {
	Daemon(daemon Daemon)

	Config(config interface{}) error

	Serve() error
	Close() error

	Restart() error
}
