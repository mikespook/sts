package rpc

import (
	"github.com/mikespook/sts/iface"
	"gopkg.in/mgo.v2/bson"
)

type Ctrl struct {
	ctrl iface.Ctrl
}

func (ctrl *Ctrl) Restart(_, _ *struct{}) (err error) {
	ctrl.ctrl.Restart()
	return
}

func (ctrl *Ctrl) Cutoff(id bson.ObjectId, reply *struct{}) error {
	return ctrl.ctrl.Cutoff(id)
}

func (ctrl *Ctrl) Kickoff(id bson.ObjectId, reply *struct{}) error {
	return ctrl.ctrl.Kickoff(id)
}
