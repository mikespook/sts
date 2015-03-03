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
	errExit1 := make(chan error, 1)
	errExit2 := make(chan error, 1)
	a.etime = time.Now()
	go func() {
		_, err := io.Copy(a.Conn, a.ch)
		errExit1 <- err
		defer close(errExit1)
	}()
	func() {
		_, err := io.Copy(a.ch, a.Conn)
		errExit2 <- err
		defer close(errExit2)
	}()
	select {
	case err = <-errExit1:
	case err = <-errExit2:
	}
	return err
}

func (a *agent) Close() error {
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
		id:   bson.NewObjectId(),
		Conn: conn,
		ch:   ch,
	}, nil
}
