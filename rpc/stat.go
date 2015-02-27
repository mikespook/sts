package rpc

import (
	"sync"
	"time"

	"github.com/mikespook/sts/model"
	"gopkg.in/mgo.v2/bson"
)

const (
	StatusConn    = "CONN"
	StatusRead    = "READ"
	StatusWritten = "WRITTEN"
)

type Stat struct {
	states model.States
}

func (stat *Stat) User(id *bson.ObjectId, reply *struct{}) error {
	return nil
}

func (stat *Stat) Conn(id bson.ObjectId, reply *struct{}) error {
	return nil
}

func (stat *Stat) Server(_, reply *struct{}) error {
	return nil
}

type Status struct {
	sync.RWMutex
	start time.Time
	data  map[string]uint64
}

func NewStatus() *Status {
	return &Status{
		start: time.Now(),
		data:  make(map[string]uint64),
	}
}

func (status *Status) Established() time.Duration {
	return time.Now().Sub(status.start)
}

func (status *Status) Inc(t string, delta uint64) uint64 {
	status.Lock()
	defer status.Unlock()
	status.data[t] += delta
	return status.data[t]
}

func (status *Status) Value(t string) uint64 {
	status.RLock()
	defer status.RUnlock()
	return status.data[t]
}
