package p2p

import (
	"encoding/gob"
	"net"

	"github.com/sirupsen/logrus"
)

type NetAddr string

func (n NetAddr) String() string  { return string(n) }
func (n NetAddr) Network() string { return "tcp" }

type Peer struct {
	conn     net.Conn
	outbound bool
}

func (p *Peer) Send(b []byte) error {
	_, err := p.conn.Write(b)
	return err
}

func (p *Peer) readLoop(msgChan chan *Message) {
	// buf := make([]byte, 1024)

	for {

		// n, err := p.conn.Read(buf)

		// if err != nil {
		// 	break
		// }

		msg := new(Message)

		if err := gob.NewDecoder(p.conn).Decode(msg); err != nil {
			logrus.Errorf("decode message error: %s", err)
			break
		}

		msgChan <- msg

		// msgChan <- &Message{
		// 	From:    p.conn.RemoteAddr(),
		// 	Payload: bytes.NewReader(buf[:n]),
		// }
	}

	p.conn.Close()
}

type TCPTransport struct {
	listenAddr string
	listener   net.Listener
	AddPeer    chan *Peer
	DeletePeer chan *Peer
}

func NewTCPTransport(addr string) *TCPTransport {
	return &TCPTransport{
		listenAddr: addr,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	ln, err := net.Listen("tcp", t.listenAddr)

	if err != nil {
		return err
	}

	t.listener = ln

	for {
		conn, err := ln.Accept()

		if err != nil {
			logrus.Error(err)
			continue
		}

		peer := &Peer{
			conn:     conn,
			outbound: false,
		}

		t.AddPeer <- peer
	}
}
