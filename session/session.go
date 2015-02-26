package session

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/mikespook/golib/log"
	"golang.org/x/crypto/ssh"
	"gopkg.in/mgo.v2/bson"
)

type session struct {
	Id bson.ObjectId

	SshConn     ssh.Conn
	ChannelChan <-chan ssh.NewChannel
	OOBReqChan  <-chan *ssh.Request
}

func (s *session) oobRequest() {
	for req := range s.OOBReqChan {
		log.Messagef("OOB Request: %+v", req)
	}
}

func (s *session) channels() {
	for newChan := range s.ChannelChan {
		ch, reqCh, err := newChan.Accept()
		if err != nil {
			log.Errorf("Channel: %s", err)
			return
		}
		chType := newChan.ChannelType()
		switch chType {
		case "session":
			go s.session(newChan, ch, reqCh)
		case "direct-tcpip":
			go s.directTcpIp(newChan, ch)
		default:
			msg := fmt.Sprintf("%s is not supported\n\r", chType)
			if _, err := ch.Write([]byte(msg)); err != nil {
				log.Errorf("Write: %s", err)
				return
			}
		}
	}
}

func (s *session) status(ch io.Writer) {
	outputs := []string{
		"\x1b[2J\x1b[1;1H",
		fmt.Sprintf("Secure Tunnel Server (%s)\r\n", s.SshConn.ServerVersion()),
		fmt.Sprintf("User: %s@%s\r\n", s.SshConn.User(),
			s.SshConn.RemoteAddr()),
		"\n* Press any key to refresh status *\r\n* Press [Ctrl+C] to disconnect *\r\n",
	}
	for _, line := range outputs {
		ch.Write([]byte(line))
	}
}

func (s *session) session(newChan ssh.NewChannel,
	ch ssh.Channel, reqChan <-chan *ssh.Request) {
	defer ch.Close()
	buf := make([]byte, 1)
LOOP:
	for {
		s.status(ch)
		if _, err := ch.Read(buf); err != nil {
			log.Errorf("Read: %s", err)
			return
		}
		switch buf[0] {
		case 0x03:
			s.Close()
			break LOOP
		default:
		}
	}
}

func parseAddr(data []byte) (addr string, err error) {
	buf := bytes.NewReader(data)
	var size uint32
	if err = binary.Read(buf, binary.BigEndian, &size); err != nil {
		return
	}
	ip := make([]byte, size)
	if err = binary.Read(buf, binary.BigEndian, ip); err != nil {
		return
	}
	var port uint32
	if err = binary.Read(buf, binary.BigEndian, &port); err != nil {
		return
	}
	addr = fmt.Sprintf("%s:%d", ip, port)
	return
}

func (s *session) directTcpIp(newChan ssh.NewChannel,
	ch ssh.Channel) {
	addr, err := parseAddr(newChan.ExtraData())
	if err != nil {
		log.Error(err)
		return
	}
	agent, err := NewAgent(addr, ch)
	if err != nil {
		log.Error(err)
		return
	}
	defer agent.Close()
	if err := agent.Serve(); err != nil {
		log.Error(err)
		return
	}
}

func (s *session) Close() error {
	return s.SshConn.Close()
}

func New(conn net.Conn, config *ssh.ServerConfig) (s *session, err error) {
	s = &session{
		Id: bson.NewObjectId(),
	}
	if s.SshConn, s.ChannelChan, s.OOBReqChan,
		err = ssh.NewServerConn(conn, config); err != nil {
		if err != io.EOF {
			return
		}
	}
	return
}

func (s *session) Serve() {
	go s.oobRequest()
	s.channels()
}

type Sessions struct {
	sync.RWMutex
	M map[bson.ObjectId]*session
}

func NewSessions() *Sessions {
	return &Sessions{
		M: make(map[bson.ObjectId]*session),
	}
}

func (m *Sessions) Add(s *session) {
	m.Lock()
	m.M[s.Id] = s
	m.Unlock()
}

func (m *Sessions) Remove(s *session) {
	m.Lock()
	delete(m.M, s.Id)
	m.Unlock()
}

func (m *Sessions) Close() {
	m.Lock()
	for k, v := range m.M {
		if err := v.Close(); err != nil {
			log.Errorf("Session Close[%s]: %s", k.Hex(), err)
		}
	}
	m.Unlock()
}
