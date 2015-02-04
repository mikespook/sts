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

type Session struct {
	SshConn     ssh.Conn
	ChannelChan <-chan ssh.NewChannel
	OOBReqChan  <-chan *ssh.Request
}

func (session *Session) OOBRequest() {
	for req := range session.OOBReqChan {
		log.Messagef("OOB Request: %+v", req)
	}
}

func (session *Session) Channels() {
	for newChan := range session.ChannelChan {
		ch, reqCh, err := newChan.Accept()
		if err != nil {
			session.Errorf("Channel: %s", err)
			return
		}
		chType := newChan.ChannelType()
		switch chType {
		case "session":
			go session.Session(newChan, ch, reqCh)
		case "direct-tcpip":
			go session.DirectTcpIp(newChan, ch)
		default:
			msg := fmt.Sprintf("%s is not supported\n\r", chType)
			if _, err := ch.Write([]byte(msg)); err != nil {
				session.Errorf("Write: %s", err)
				return
			}
		}
	}
}

func (session *Session) status(ch io.Writer) {
	outputs := []string{
		"\x1b[2J\x1b[1;1H",
		fmt.Sprintf("Secure Tunnel Server (%s)\r\n", session.SshConn.ServerVersion()),
		fmt.Sprintf("User: %s@%s\r\n", session.SshConn.User(),
			session.SshConn.RemoteAddr()),
		"\n* Press any key to refresh status *\r\n* Press [Ctrl+C] to disconnect *\r\n",
	}
	for _, line := range outputs {
		ch.Write([]byte(line))
	}
}

func (session *Session) Session(newChan ssh.NewChannel,
	ch ssh.Channel, reqChan <-chan *ssh.Request) {
	defer ch.Close()
	buf := make([]byte, 1)
LOOP:
	for {
		session.status(ch)
		if _, err := ch.Read(buf); err != nil {
			session.Errorf("Read: %s", err)
			return
		}
		switch buf[0] {
		case 0x03:
			session.Close()
			break LOOP
		default:
		}
	}
}

func (session *Session) Errorf(format string, err error) {
	log.Errorf("%s [%s]", fmt.Sprintf(format, err), session.SshConn.User())
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

func (session *Session) DirectTcpIp(newChan ssh.NewChannel,
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

func (session *Session) Close() error {
	return session.SshConn.Close()
}

func doSession(conn net.Conn, config *ssh.ServerConfig) {
	defer func() {
		conn.Close()
		log.Messagef("Disconnect: %s", conn.RemoteAddr())
	}()
	log.Messagef("Connect: %s", conn.RemoteAddr())
	session := &Session{}
	var err error
	if session.SshConn, session.ChannelChan, session.OOBReqChan,
		err = ssh.NewServerConn(conn, config); err != nil {
		if err != io.EOF {
			log.Errorf("SSH-Connect: %s", err)
		}
		return
	}
	log.Messagef("SSH-Connect: %s@%s (%s)", session.SshConn.User(),
		session.SshConn.RemoteAddr(), session.SshConn.ClientVersion())
	go session.OOBRequest()
	session.Channels()
	log.Messagef("SSH-Disconnect: %s@%s (%s)", session.SshConn.User(),
		session.SshConn.RemoteAddr(), session.SshConn.ClientVersion())
	return
}
