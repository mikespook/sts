package tunnel

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/mikespook/golib/log"
	"github.com/mikespook/sts/model"
	"golang.org/x/crypto/ssh"
	"gopkg.in/mgo.v2/bson"
)

type session struct {
	id bson.ObjectId

	sshConn  ssh.Conn
	channels <-chan ssh.NewChannel
	oobReqs  <-chan *ssh.Request

	states model.States
}

func newSession(conn net.Conn, config *ssh.ServerConfig, states model.States) (s *session, err error) {
	s = &session{
		id:     bson.NewObjectId(),
		states: states,
	}
	if s.sshConn, s.channels, s.oobReqs,
		err = ssh.NewServerConn(conn, config); err != nil {
		if err != io.EOF {
			return
		}
	}
	return
}

func (s *session) Id() bson.ObjectId {
	return s.id
}

func (s *session) Metadata() ssh.ConnMetadata {
	return s.sshConn
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
		fmt.Sprintf("Secure Tunnel Server (%s)\r\n", s.sshConn.ServerVersion()),
		fmt.Sprintf("User: %s@%s\r\n", s.sshConn.User(),
			s.sshConn.RemoteAddr()),
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
	a, err := newAgent(addr, ch)
	if err != nil {
		log.Error(err)
		return
	}
	defer a.Close()
	s.states.Agents().Add(a)
	defer s.states.Agents().Remove(a)
	if err := a.Serve(); err != nil {
		log.Error(err)
		return
	}
}

func (s *session) Close() error {
	return s.sshConn.Close()
}

func (s *session) Serve() {
	go s.serveOOBRequest()
	log.Messagef("SSH-Connect: %s [%s@%s] (%s)", s.id.Hex(), s.sshConn.User(),
		s.sshConn.RemoteAddr(), s.sshConn.ClientVersion())
	s.serveChannels()
	log.Messagef("SSH-Disconnect: %s", s.id.Hex())
}
