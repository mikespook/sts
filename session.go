package sts

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/mikespook/golib/log"
	"golang.org/x/crypto/ssh"
)

type session struct {
	SshConn     ssh.Conn
	ChannelChan <-chan ssh.NewChannel
	OOBReqChan  <-chan *ssh.Request
}

func (s *session) OOBRequest() {
	for req := range s.OOBReqChan {
		log.Messagef("OOB Request: %+v", req)
	}
}

func (s *session) Channels() {
	for newChan := range s.ChannelChan {
		ch, reqCh, err := newChan.Accept()
		if err != nil {
			s.Errorf("Channel: %s", err)
			return
		}
		chType := newChan.ChannelType()
		switch chType {
		case "session":
			go s.Session(newChan, ch, reqCh)
		case "direct-tcpip":
			go s.DirectTcpIp(newChan, ch)
		default:
			msg := fmt.Sprintf("%s is not supported\n\r", chType)
			if _, err := ch.Write([]byte(msg)); err != nil {
				s.Errorf("Write: %s", err)
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

func (s *session) Session(newChan ssh.NewChannel,
	ch ssh.Channel, reqChan <-chan *ssh.Request) {
	defer ch.Close()
	buf := make([]byte, 1)
LOOP:
	for {
		s.status(ch)
		if _, err := ch.Read(buf); err != nil {
			s.Errorf("Read: %s", err)
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

func (s *session) Errorf(format string, err error) {
	log.Errorf("%s [%s]", fmt.Sprintf(format, err), s.SshConn.User())
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

func (s *session) DirectTcpIp(newChan ssh.NewChannel,
	ch ssh.Channel) {
	defer ch.Close()
	addr, err := parseAddr(newChan.ExtraData())
	if err != nil {
		log.Error(err)
		return
	}
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Error(err)
		return
	}
	defer conn.Close()
	go io.Copy(conn, ch)
	io.Copy(ch, conn)
}

func (s *session) Close() error {
	return s.SshConn.Close()
}

func doSession(conn net.Conn, config *ssh.ServerConfig) {
	defer func() {
		conn.Close()
		log.Messagef("Disconnect: %s", conn.RemoteAddr())
	}()
	log.Messagef("Connect: %s", conn.RemoteAddr())
	s := &session{}
	var err error
	if s.SshConn, s.ChannelChan, s.OOBReqChan,
		err = ssh.NewServerConn(conn, config); err != nil {
		if err != io.EOF {
			log.Errorf("SSH-Connect: %s", err)
		}
		return
	}
	log.Messagef("SSH-Connect: %s@%s (%s)", s.SshConn.User(),
		s.SshConn.RemoteAddr(), s.SshConn.ClientVersion())
	go s.OOBRequest()
	s.Channels()
	log.Messagef("SSH-Disconnect: %s@%s (%s)", s.SshConn.User(),
		s.SshConn.RemoteAddr(), s.SshConn.ClientVersion())
	return
}
