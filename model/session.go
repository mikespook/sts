package model

import (
	"sync"
	"time"

	"github.com/mikespook/sts/iface"
	"gopkg.in/mgo.v2/bson"
)

const (
	StatSession = "session"
)

type Session struct {
	Agents        map[bson.ObjectId]*Agent
	ETime         time.Time
	ClientVersion []byte
	Id            bson.ObjectId
	LocalAddr     string
	RemoteAddr    string
	ServerVersion []byte
	User          string
}

type Sessions struct {
	M map[bson.ObjectId]*Session
}

type SessionIfaces struct {
	sync.RWMutex
	M map[bson.ObjectId]iface.Session
}

func NewSessionIfaces() *SessionIfaces {
	return &SessionIfaces{
		M: make(map[bson.ObjectId]iface.Session),
	}
}

func (m *SessionIfaces) Add(s iface.Session) {
	m.Lock()
	m.M[s.Id()] = s
	m.Unlock()
}

func (m *SessionIfaces) Remove(s iface.Session) {
	m.Lock()
	delete(m.M, s.Id())
	m.Unlock()
}
