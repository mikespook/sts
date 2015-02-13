package sts

import (
	"github.com/mikespook/golib/log"
	"gopkg.in/mgo.v2/bson"
)

type Ctrl struct {
	server *Server
}

func (ctrl *Ctrl) Restart(_, reply *struct{}) (err error) {
	log.Messagef("Restarting: addr=%s, keys=%+v, pwd=%s",
		ctrl.server.config.Addr, ctrl.server.config.Keys,
		ctrl.server.config.Pwd)

	ctrl.server.reboot()
	return
}

func (ctrl *Ctrl) Cutoff(id bson.ObjectId, reply *struct{}) error {
	return ctrl.server.cutoff(id)
}

func (ctrl *Ctrl) Kickoff(id bson.ObjectId, reply *struct{}) error {
	return ctrl.server.kickoff(id)
}
