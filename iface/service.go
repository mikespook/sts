package iface

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Keeper interface {
	AddAgent(a Agent)
	RemoveAgent(a Agent)
	Agent(id bson.ObjectId) Agent
	Agents() map[bson.ObjectId]Agent

	AddSession(s Session)
	RemoveSession(s Session)
	Session(id bson.ObjectId) Session
	Sessions() map[bson.ObjectId]Session

	Restart()
	Cutoff(id bson.ObjectId)
	Kickoff(id bson.ObjectId)

	Stat() Stat
}

type Service interface {
	Config(config interface{}) error

	Serve() error
	Close() error

	ETime() time.Time
}
