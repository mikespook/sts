package rpc

import (
	"github.com/mikespook/sts/bus"
	"gopkg.in/mgo.v2/bson"
)

type Ctrl struct {
	bus bus.Ctrl
}

func (ctrl *Ctrl) Restart(_, _ *struct{}) (err error) {
	ctrl.bus.Restart()
	return
}

func (ctrl *Ctrl) Cutoff(id bson.ObjectId, reply *struct{}) error {
	return ctrl.bus.Cutoff(id)
}

func (ctrl *Ctrl) Kickoff(id bson.ObjectId, reply *struct{}) error {
	return ctrl.bus.Kickoff(id)
}
