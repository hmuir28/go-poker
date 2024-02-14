package p2p

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"sync"
)

type Peer struct {
	conn net.Conn
}

func (p *Peer) Send(b []byte) error {
	_, err := p.conn.Write(b)
	return err
}

type ServerConfig struct {
	ListenAddr string
	Version    string
}

type Message struct {
	Payload io.Reader
	From    net.Addr
}

type Server struct {
	ServerConfig

	handler    Handler
	listener   net.Listener
	mu         sync.RWMutex
	peers      map[net.Addr]*Peer
	addPeer    chan *Peer
	deletePeer chan *Peer
	msgChan    chan *Message
}

func NewServer(config ServerConfig) *Server {
	return &Server{
		handler:      &DefaultHandler{},
		ServerConfig: config,
		peers:        make(map[net.Addr]*Peer),
		addPeer:      make(chan *Peer),
		msgChan:      make(chan *Message),
	}
}

func (s *Server) Start() {
	go s.loop()

	if err := s.listen(); err != nil {
		panic(err)
	}

	fmt.Printf("game server running on port %s\n", s.ListenAddr)

	s.acceptLoop()
}

func (s *Server) handleConn(p *Peer) {
	buf := make([]byte, 1024)

	for {

		n, err := p.conn.Read(buf)

		if err != nil {
			break
		}

		s.msgChan <- &Message{
			From:    p.conn.RemoteAddr(),
			Payload: bytes.NewReader(buf[:n]),
		}
	}

	s.deletePeer <- p
}

func (s *Server) Connect(addr string) error {
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		return err
	}

	peer := &Peer{
		conn: conn,
	}

	s.addPeer <- peer

	return peer.Send([]byte(s.Version))
}

func (s *Server) acceptLoop() {
	for {

		conn, err := s.listener.Accept()

		if err != nil {
			panic(err)
		}

		peer := &Peer{
			conn: conn,
		}

		s.addPeer <- peer

		peer.Send([]byte(s.Version))

		go s.handleConn(peer)
	}
}

func (s *Server) listen() error {
	ln, err := net.Listen("tcp", s.ListenAddr)

	if err != nil {
		return err
	}

	s.listener = ln

	return nil
}

func (s *Server) Stop() {

}

func (s *Server) loop() {
	for {
		select {

		case peer := <-s.deletePeer:
			delete(s.peers, peer.conn.RemoteAddr())
			fmt.Printf("player disconnected %s \n", peer.conn.RemoteAddr())

		case peer := <-s.addPeer:
			s.peers[peer.conn.RemoteAddr()] = peer
			fmt.Printf("new player connected %s\n", peer.conn.RemoteAddr())

		case msg := <-s.msgChan:
			if err := s.handler.HandleMessage(msg); err != nil {
				panic(err)
			}
		}
	}
}
