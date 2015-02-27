package model

import (
	"net"
	"sync"

	"golang.org/x/crypto/ssh"
	"gopkg.in/mgo.v2/bson"
)

type States interface {
	Sessions() *Sessions
	Agents() *Agents
}

type Sessions struct {
	sync.RWMutex
	M map[bson.ObjectId]Session
}

type Session interface {
	Id() bson.ObjectId
	Metadata() ssh.ConnMetadata
	Close() error
}

func NewSessions() *Sessions {
	return &Sessions{
		M: make(map[bson.ObjectId]Session),
	}
}

func (m *Sessions) Add(s Session) {
	m.Lock()
	m.M[s.Id()] = s
	m.Unlock()
}

func (m *Sessions) Remove(s Session) {
	m.Lock()
	delete(m.M, s.Id())
	m.Unlock()
}

type Agents struct {
	sync.RWMutex
	M map[bson.ObjectId]Agent
}

type Agent interface {
	Id() bson.ObjectId
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close() error
}

func NewAgents() *Agents {
	return &Agents{
		M: make(map[bson.ObjectId]Agent),
	}
}

func (m *Agents) Add(a Agent) {
	m.Lock()
	m.M[a.Id()] = a
	m.Unlock()
}

func (m *Agents) Remove(a Agent) {
	m.Lock()
	delete(m.M, a.Id())
	m.Unlock()
}
