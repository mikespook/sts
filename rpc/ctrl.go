package rpc

import (
	"github.com/mikespook/sts/model"
	"gopkg.in/mgo.v2/bson"
)

type Ctrl struct {
	keeper model.Keeper
}

func (ctrl *Ctrl) Restart(_, _ *struct{}) error {
	ctrl.keeper.Restart()
	return nil
}

func (ctrl *Ctrl) Cutoff(id bson.ObjectId, reply *struct{}) error {
	ctrl.keeper.Cutoff(id)
	return nil
}

func (ctrl *Ctrl) Kickoff(id bson.ObjectId, reply *struct{}) error {
	ctrl.keeper.Kickoff(id)
	return nil
}
