package p2p

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
)

type GameVariant uint8

func (gv GameVariant) String() string {
	switch gv {
	case TEXAS_HOLDEM:
		return "TEXAS HOLDEM"
	case Other:
		return "Other"
	default:
		return "unknown"
	}
}

const (
	TEXAS_HOLDEM GameVariant = iota
	Other
)

type ServerConfig struct {
	ListenAddr  string
	Version     string
	GameVariant GameVariant
}

type Server struct {
	ServerConfig

	transport  *TCPTransport
	listener   net.Listener
	peers      map[net.Addr]*Peer
	addPeer    chan *Peer
	deletePeer chan *Peer
	msgChan    chan *Message

	gameState *GameState
}

func NewServer(config ServerConfig) *Server {
	s := &Server{
		ServerConfig: config,
		peers:        make(map[net.Addr]*Peer),
		addPeer:      make(chan *Peer),
		msgChan:      make(chan *Message),
		gameState:    NewGameState(),
	}

	tr := NewTCPTransport(s.ListenAddr)
	s.transport = tr

	tr.AddPeer = s.addPeer
	tr.DeletePeer = s.deletePeer

	return s
}

func (s *Server) Start() {
	go s.loop()

	logrus.WithFields(logrus.Fields{
		"port":    s.ListenAddr,
		"variant": s.GameVariant,
	}).Info("Started new game server")
	fmt.Printf("game server running on port %s\n", s.ListenAddr)

	s.transport.ListenAndAccept()
}

func (s *Server) SendHandshake(p *Peer) error {
	hs := &Handshake{
		GameVariant: s.GameVariant,
		Version:     s.Version,
	}

	buf := new(bytes.Buffer)

	if err := gob.NewEncoder(buf).Encode(hs); err != nil {
		return err
	}

	return p.Send(buf.Bytes())
}

func (s *Server) Connect(addr string) error {
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		return err
	}

	peer := &Peer{
		conn:     conn,
		outbound: true,
	}

	s.addPeer <- peer

	return s.SendHandshake(peer)
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
			logrus.WithFields(logrus.Fields{
				"addr": peer.conn.RemoteAddr(),
			}).Info("New player disconnected")

			delete(s.peers, peer.conn.RemoteAddr())

		case peer := <-s.addPeer:
			if err := s.handshake(peer); err != nil {
				logrus.Errorf("handshake with incoming player failed: %s", err)
				peer.conn.Close()

				delete(s.peers, peer.conn.RemoteAddr())
				continue
			}

			go peer.readLoop(s.msgChan)

			if !peer.outbound {
				if err := s.SendHandshake(peer); err != nil {
					logrus.Errorf("failed to send handshake with peer: %s", err)
					peer.conn.Close()
					delete(s.peers, peer.conn.RemoteAddr())
					continue
				}
			}

			logrus.WithFields(logrus.Fields{
				"addr": peer.conn.RemoteAddr(),
			}).Info("handshake succesfull: player connected")

			s.peers[peer.conn.RemoteAddr()] = peer

		case msg := <-s.msgChan:
			if err := s.handleMessage(msg); err != nil {
				panic(err)
			}
		}
	}
}

type Handshake struct {
	Version     string
	GameVariant GameVariant
}

func (s *Server) handshake(p *Peer) error {
	hs := &Handshake{}

	if err := gob.NewDecoder(p.conn).Decode(hs); err != nil {
		return err
	}

	if s.GameVariant != hs.GameVariant {
		return fmt.Errorf("game variant does not match %s\n", hs.GameVariant)
	}

	if s.Version != hs.Version {
		return fmt.Errorf("invalid version %s\n", hs.Version)
	}

	logrus.WithFields(logrus.Fields{
		"peer":    p.conn.RemoteAddr(),
		"version": hs.Version,
		"variant": hs.GameVariant,
	}).Info("received handshake")

	return nil
}

func (s *Server) handleMessage(msg *Message) error {
	fmt.Printf("%+v \n", msg)
	return nil
}

func Init() {
	gob.Register(Handshake{})
}
