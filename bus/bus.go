package bus

import (
	"net"

	"golang.org/x/crypto/ssh"
	"gopkg.in/mgo.v2/bson"
)

type Bus interface {
	Errorf(format string, args ...interface{})
	Messagef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Warningf(format string, args ...interface{})

	Session(conn net.Conn, config *ssh.ServerConfig)
	Close()

	Ctrl() Ctrl
	Stat() Stat
}

type Service interface {
	Bus(bus Bus)

	Config(config interface{}) error

	Serve() error
	Close() error

	Restart() error
}

type Ctrl interface {
	Restart() error
	Cutoff(id bson.ObjectId) error
	Kickoff(id bson.ObjectId) error
}

type Stat interface {
}
