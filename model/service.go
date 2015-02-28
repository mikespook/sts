package model

import "gopkg.in/mgo.v2/bson"

type Keeper interface {
	AddAgent(a Agent)
	RemoveAgent(a Agent)
	AddSession(s Session)
	RemoveSession(s Session)

	Restart()
	Cutoff(id bson.ObjectId)
	Kickoff(id bson.ObjectId)
}

type Service interface {
	Config(config interface{}) error

	Serve() error
	Close() error

	Restart() error
}
