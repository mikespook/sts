package iface

type Daemon interface {
	Close()

	Ctrl() Ctrl
	Stat() Stat
}
