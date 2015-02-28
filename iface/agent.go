package iface

import (
	"net"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Agent interface {
	Id() bson.ObjectId
	User() string
	SessionId() bson.ObjectId
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close() error
	ETime() time.Time
}
