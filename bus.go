package sts

import (
	"time"

	"github.com/mikespook/sts/iface"
	"github.com/mikespook/sts/model"
	"gopkg.in/mgo.v2/bson"
)

type bus struct {
	sts         *Sts
	sessions    *model.SessionIfaces
	agents      *model.AgentIfaces
	established time.Time
}

func newBus(sts *Sts) *bus {
	return &bus{
		sts:      sts,
		sessions: model.NewSessionIfaces(),
		agents:   model.NewAgentIfaces(),
	}
}

func (k *bus) Stat() iface.Stat {
	return k
}

func (k *bus) AddSession(s iface.Session) {
	k.sessions.Add(s)
}

func (k *bus) RemoveSession(s iface.Session) {
	k.sessions.Remove(s)
}

func (k *bus) Session(id bson.ObjectId) iface.Session {
	return k.sessions.M[id]
}

func (k *bus) Sessions() map[bson.ObjectId]iface.Session {
	return k.sessions.M
}

func (k *bus) AddAgent(a iface.Agent) {
	k.agents.Add(a)
}

func (k *bus) RemoveAgent(a iface.Agent) {
	k.agents.Remove(a)
}

func (k *bus) Agent(id bson.ObjectId) iface.Agent {
	return k.agents.M[id]
}

func (k *bus) Agents() map[bson.ObjectId]iface.Agent {
	return k.agents.M
}

func (k *bus) Restart() {
	k.sts.restart()
}

func (k *bus) Cutoff(id bson.ObjectId) {
	k.agents.Lock()
	if a, ok := k.agents.M[id]; ok {
		a.Close()
		delete(k.agents.M, id)
	}
	k.agents.Unlock()
}

func (k *bus) Kickoff(id bson.ObjectId) {
	k.sessions.Lock()
	if s, ok := k.sessions.M[id]; ok {
		s.Close()
		delete(k.sessions.M, id)
	}
	k.sessions.Unlock()
}

func (k *bus) ETime() time.Time {
	return k.sts.services[Tunnel].ETime()
}

func (k *bus) Aggregate(key string) int {
	switch key {
	case model.StatSession:
		return len(k.sessions.M)
	case model.StatAgent:
		return len(k.agents.M)
	}
	return 0
}
