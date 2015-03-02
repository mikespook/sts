package tunnel

import (
	"io"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
	"gopkg.in/mgo.v2/bson"
)

type agent struct {
	id bson.ObjectId
	ch ssh.Channel
	net.Conn

	session *session
	etime   time.Time
	errExit chan error
}

func (a *agent) Id() bson.ObjectId {
	return a.id
}

func (a *agent) SessionId() bson.ObjectId {
	return a.session.id
}

func (a *agent) User() string {
	return a.session.User()
}

func (a *agent) Serve() (err error) {
	a.etime = time.Now()
	go func() {
		_, err := io.Copy(a.Conn, a.ch)
		a.errExit <- err
	}()
	go func() {
		_, err := io.Copy(a.ch, a.Conn)
		a.errExit <- err
	}()
	for err = range a.errExit {
		break
	}
	return err
}

func (a *agent) Close() error {
	close(a.errExit)
	return a.Conn.Close()
}

func (a *agent) ETime() time.Time {
	return a.etime
}

func newAgent(addr string, ch ssh.Channel) (*agent, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &agent{
		id:      bson.NewObjectId(),
		Conn:    conn,
		ch:      ch,
		errExit: make(chan error),
	}, nil
}
