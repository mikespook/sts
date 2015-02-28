package iface

import (
	"time"

	"golang.org/x/crypto/ssh"
	"gopkg.in/mgo.v2/bson"
)

type Session interface {
	Agents() map[bson.ObjectId]Agent
	Id() bson.ObjectId
	ssh.ConnMetadata
	Close() error
	ETime() time.Time
}
