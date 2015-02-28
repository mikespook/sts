package sts

import (
	"time"

	"github.com/mikespook/sts/iface"
	"github.com/mikespook/sts/model"
	"gopkg.in/mgo.v2/bson"
)

type keeper struct {
	sts         *Sts
	sessions    *model.SessionIfaces
	agents      *model.AgentIfaces
	established time.Time
}

func newKeeper(sts *Sts) *keeper {
	return &keeper{
		sts:      sts,
		sessions: model.NewSessionIfaces(),
		agents:   model.NewAgentIfaces(),
	}
}

func (k *keeper) Stat() iface.Stat {
	return k
}

func (k *keeper) AddSession(s iface.Session) {
	k.sessions.Add(s)
}

func (k *keeper) RemoveSession(s iface.Session) {
	k.sessions.Remove(s)
}

func (k *keeper) Session(id bson.ObjectId) iface.Session {
	return k.sessions.M[id]
}

func (k *keeper) Sessions() map[bson.ObjectId]iface.Session {
	return k.sessions.M
}

func (k *keeper) AddAgent(a iface.Agent) {
	k.agents.Add(a)
}

func (k *keeper) RemoveAgent(a iface.Agent) {
	k.agents.Remove(a)
}

func (k *keeper) Agent(id bson.ObjectId) iface.Agent {
	return k.agents.M[id]
}

func (k *keeper) Agents() map[bson.ObjectId]iface.Agent {
	return k.agents.M
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

func (k *keeper) ETime() time.Time {
	return k.sts.services[Tunnel].ETime()
}

func (k *keeper) Aggregate(key string) int {
	switch key {
	case model.StatSession:
		return len(k.sessions.M)
	case model.StatAgent:
		return len(k.agents.M)
	}
	return 0
}
