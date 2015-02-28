package rpc

import (
	"github.com/mikespook/sts/iface"
	"gopkg.in/mgo.v2/bson"
)

type rpcCtrl struct {
	keeper iface.Keeper
}

func (ctrl *rpcCtrl) Restart(_, _ *struct{}) error {
	ctrl.keeper.Restart()
	return nil
}

func (ctrl *rpcCtrl) Cutoff(id bson.ObjectId, _ *struct{}) error {
	ctrl.keeper.Cutoff(id)
	return nil
}

func (ctrl *rpcCtrl) Kickoff(id bson.ObjectId, _ *struct{}) error {
	ctrl.keeper.Kickoff(id)
	return nil
}
