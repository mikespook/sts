package rpc

import (
	"github.com/mikespook/sts/iface"
	"gopkg.in/mgo.v2/bson"
)

type rpcCtrl struct {
	bus iface.Bus
}

func (ctrl *rpcCtrl) Restart(_ interface{}, _ *struct{}) error {
	ctrl.bus.Restart()
	return nil
}

func (ctrl *rpcCtrl) Cutoff(id bson.ObjectId, _ *struct{}) error {
	ctrl.bus.Cutoff(id)
	return nil
}

func (ctrl *rpcCtrl) Kickoff(id bson.ObjectId, _ *struct{}) error {
	ctrl.bus.Kickoff(id)
	return nil
}
