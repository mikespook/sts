package tunnel

import (
	"io"
	"net"

	"golang.org/x/crypto/ssh"
	"gopkg.in/mgo.v2/bson"
)

type agent struct {
	id bson.ObjectId
	ch ssh.Channel
	net.Conn
}

func (a *agent) Id() bson.ObjectId {
	return a.id
}

func (a *agent) Serve() error {
	defer a.ch.Close()
	go io.Copy(a.Conn, a.ch)
	_, err := io.Copy(a.ch, a.Conn)
	return err
}

func (a *agent) Close() error {
	return a.Conn.Close()
}

func newAgent(addr string, ch ssh.Channel) (*agent, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &agent{
		id:   bson.NewObjectId(),
		Conn: conn,
		ch:   ch,
	}, nil
}
