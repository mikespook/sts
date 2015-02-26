package iface

import "gopkg.in/mgo.v2/bson"

type Stat interface {
	All() map[bson.ObjectId]Item
	Add(item Item)
	Remove(item Item)
}

type Item interface {
	Id() bson.ObjectId
}
