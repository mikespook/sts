package model

type Service interface {
	Config(config interface{}) error

	Serve() error
	Close() error

	Restart() error
}
