package rpc

import (
	"github.com/mikespook/golib/log"
	"github.com/mikespook/sts/iface"
	"gopkg.in/mgo.v2/bson"
)

type rpcCtrl struct {
	bus iface.Bus
}

func (ctrl *rpcCtrl) Restart(_ struct{}, _ *struct{}) error {
	log.Message("RPC: Restart")
	ctrl.bus.Restart()
	return nil
}

func (ctrl *rpcCtrl) Cutoff(id bson.ObjectId, _ *struct{}) error {
	log.Messagef("RPC: Cut off %s", id.Hex())
	ctrl.bus.Cutoff(id)
	return nil
}

func (ctrl *rpcCtrl) Kickoff(id bson.ObjectId, _ *struct{}) error {
	log.Messagef("RPC: Kick off %s", id.Hex())
	ctrl.bus.Kickoff(id)
	return nil
}
