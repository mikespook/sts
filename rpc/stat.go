package rpc

import (
	"fmt"

	"github.com/mikespook/sts/iface"
	"github.com/mikespook/sts/model"
	"gopkg.in/mgo.v2/bson"
)

const (
	StatusConn    = "CONN"
	StatusRead    = "READ"
	StatusWritten = "WRITTEN"
)

type rpcStat struct {
	bus iface.Bus
}

func sessionConv(i iface.Session, m *model.Session) {
	ia := i.Agents()
	m.Agents = make(map[bson.ObjectId]*model.Agent)
	for k, v := range ia {
		var ma model.Agent
		agentConv(v, &ma)
		m.Agents[k] = &ma
	}
	m.ETime = i.ETime()
	m.ClientVersion = i.ClientVersion()
	m.Id = i.Id()
	m.LocalAddr = i.LocalAddr().String()
	m.RemoteAddr = i.RemoteAddr().String()
	m.ServerVersion = i.ServerVersion()
	m.User = i.User()
}

func agentConv(i iface.Agent, m *model.Agent) {
	m.ETime = i.ETime()
	m.Id = i.Id()
	m.LocalAddr = i.LocalAddr().String()
	m.RemoteAddr = i.RemoteAddr().String()
	m.SessionId = i.SessionId()
	m.User = i.User()
}

func (stat *rpcStat) Sessions(user string, s *model.Sessions) error {
	s.M = make(map[bson.ObjectId]*model.Session)
	i := stat.bus.Sessions()
	for k, v := range i {
		if user == "" || v.User() == user {
			var m model.Session
			sessionConv(v, &m)
			s.M[k] = &m
		}
	}
	return nil
}

func (stat *rpcStat) Agents(user string, a *model.Agents) error {
	a.M = make(map[bson.ObjectId]*model.Agent)
	i := stat.bus.Agents()
	for k, v := range i {
		if user == "" || v.User() == user {
			var m model.Agent
			agentConv(v, &m)
			a.M[k] = &m
		}
	}
	return nil
}

func (stat *rpcStat) Stat(_ struct{}, s *model.Stat) error {
	i := stat.bus.Stat()
	s.Sessions = i.Aggregate(model.StatSession)
	s.Agents = i.Aggregate(model.StatAgent)
	s.ETime = i.ETime()
	return nil
}

func (stat *rpcStat) Session(id bson.ObjectId, s *model.Session) error {
	i := stat.bus.Session(id)
	if i == nil {
		return fmt.Errorf("Session %s not found", id.Hex())
	}
	sessionConv(i, s)
	return nil
}

func (stat *rpcStat) Agent(id bson.ObjectId, a *model.Agent) error {
	i := stat.bus.Agent(id)
	if i == nil {
		return fmt.Errorf("Agent %s not found", id.Hex())
	}
	agentConv(i, a)
	return nil
}
