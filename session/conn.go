package session

import (
	"io"
	"net"
	"sync"

	"github.com/mikespook/golib/log"
	"golang.org/x/crypto/ssh"
	"gopkg.in/mgo.v2/bson"
)

type agent struct {
	Id bson.ObjectId
	Ch ssh.Channel
	net.Conn
}

func (a *agent) Serve() error {
	defer a.Conn.Close()
	go io.Copy(a.Conn, a.Ch)
	_, err := io.Copy(a.Ch, a.Conn)
	return err
}

func (a *agent) Close() error {
	return a.Ch.Close()
}

func NewAgent(addr string, ch ssh.Channel) (*agent, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &agent{
		Id:   bson.NewObjectId(),
		Conn: conn,
		Ch:   ch,
	}, nil
}

type Agents struct {
	sync.RWMutex
	M map[bson.ObjectId]*agent
}

func NewAgents() *Agents {
	return &Agents{
		M: make(map[bson.ObjectId]*agent),
	}
}

func (m *Agents) Add(a *agent) {
	m.Lock()
	m.M[a.Id] = a
	m.Unlock()
}

func (m *Agents) Remove(a *agent) {
	m.Lock()
	delete(m.M, a.Id)
	m.Unlock()
}

func (m *Agents) Close() {
	m.Lock()
	for k, v := range m.M {
		if err := v.Close(); err != nil {
			log.Errorf("Agent Close[%s]: %s", k.Hex(), err)
		}
	}
	m.Unlock()
}
