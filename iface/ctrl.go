package iface

import "gopkg.in/mgo.v2/bson"

type Ctrl interface {
	Restart() error
	Cutoff(id bson.ObjectId) error
	Kickoff(id bson.ObjectId) error
}
