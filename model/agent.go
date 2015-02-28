package model

import (
	"encoding/gob"
	"net"
	"sync"
	"time"

	"github.com/mikespook/sts/iface"
	"gopkg.in/mgo.v2/bson"
)

func init() {
	gob.Register(net.TCPAddr{})
}

const (
	StatAgent = "agent"
)

type Agent struct {
	ETime      time.Time
	Id         bson.ObjectId
	LocalAddr  string
	RemoteAddr string
	SessionId  bson.ObjectId
	User       string
}

type Agents struct {
	M map[bson.ObjectId]*Agent
}

type AgentIfaces struct {
	sync.RWMutex
	M map[bson.ObjectId]iface.Agent
}

func NewAgentIfaces() *AgentIfaces {
	return &AgentIfaces{
		M: make(map[bson.ObjectId]iface.Agent),
	}
}

func (m *AgentIfaces) Add(a iface.Agent) {
	m.Lock()
	m.M[a.Id()] = a
	m.Unlock()
}

func (m *AgentIfaces) Remove(a iface.Agent) {
	m.Lock()
	delete(m.M, a.Id())
	m.Unlock()
}
