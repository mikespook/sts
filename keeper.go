package sts

import (
	"github.com/mikespook/sts/model"
	"gopkg.in/mgo.v2/bson"
)

type keeper struct {
	sts      *Sts
	sessions *model.Sessions
	agents   *model.Agents
}

func newKeeper(sts *Sts) *keeper {
	return &keeper{
		sts:      sts,
		sessions: model.NewSessions(),
		agents:   model.NewAgents(),
	}
}

func (k *keeper) AddSession(s model.Session) {
	k.sessions.Add(s)
}

func (k *keeper) RemoveSession(s model.Session) {
	k.sessions.Add(s)
}

func (k *keeper) AddAgent(a model.Agent) {
	k.agents.Add(a)
}

func (k *keeper) RemoveAgent(a model.Agent) {
	k.agents.Add(a)
}

func (k *keeper) Restart() {
	k.sts.restart()
}

func (k *keeper) Cutoff(id bson.ObjectId) {
	k.agents.Lock()
	k.agents.M[id].Close()
	delete(k.agents.M, id)
	k.agents.Unlock()
}

func (k *keeper) Kickoff(id bson.ObjectId) {
	k.sessions.Lock()
	k.sessions.M[id].Close()
	delete(k.sessions.M, id)
	k.sessions.Unlock()
}
