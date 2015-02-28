package tunnel

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/mikespook/golib/log"
	"github.com/mikespook/sts/iface"
	"golang.org/x/crypto/ssh"
	"gopkg.in/mgo.v2/bson"
)

type session struct {
	id bson.ObjectId

	ssh.Conn
	channels <-chan ssh.NewChannel
	oobReqs  <-chan *ssh.Request

	keeper iface.Keeper
	etime  time.Time
}

func newSession(conn net.Conn, config *ssh.ServerConfig) (s *session, err error) {
	s = &session{
		id: bson.NewObjectId(),
	}
	if s.Conn, s.channels, s.oobReqs,
		err = ssh.NewServerConn(conn, config); err != nil {
		if err != io.EOF {
			return
		}
	}
	return
}

func (s *session) ETime() time.Time {
	return s.etime
}

func (s *session) Id() bson.ObjectId {
	return s.id
}

func (s *session) serveOOBRequest() {
	for req := range s.oobReqs {
		log.Messagef("OOB Request: %+v", req)
	}
}

func (s *session) serveChannels() {
	for newChan := range s.channels {
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
		fmt.Sprintf("Secure Tunnel Server (%s)\r\n", s.ServerVersion()),
		fmt.Sprintf("User: %s@%s\r\n", s.User(),
			s.RemoteAddr()),
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
	defer ch.Close()
	addr, err := parseAddr(newChan.ExtraData())
	if err != nil {
		log.Error(err)
		return
	}
	a, err := newAgent(addr, ch)
	if err != nil {
		log.Error(err)
		return
	}
	a.session = s
	defer a.Close()
	s.keeper.AddAgent(a)
	defer s.keeper.RemoveAgent(a)
	if err := a.Serve(); err != nil {
		log.Error(err)
		return
	}
}

func (s *session) Close() error {
	return s.Conn.Close()
}

func (s *session) Agents() map[bson.ObjectId]iface.Agent {
	agents := make(map[bson.ObjectId]iface.Agent)
	all := s.keeper.Agents()
	for k, v := range all {
		if v.User() == s.User() {
			agents[k] = v
		}
	}
	return agents
}

func (s *session) Serve() {
	s.etime = time.Now()
	go s.serveOOBRequest()
	log.Messagef("SSH-Connect: %s [%s@%s] (%s)", s.id.Hex(), s.User(),
		s.RemoteAddr(), s.ClientVersion())
	s.serveChannels()
	log.Messagef("SSH-Disconnect: %s", s.id.Hex())
}
