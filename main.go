package main

import (
	"time"

	"github.com/hmuir28/go-poker/p2p"
)

func makeServer(addr string) *p2p.Server {
	config := p2p.ServerConfig{
		Version:     "GO POKER v0.1-alpha",
		ListenAddr:  ":3000",
		GameVariant: p2p.TEXAS_HOLDEM,
	}
	server := p2p.NewServer(config)

	go server.Start()

	time.Sleep(1 * time.Second)

	return server
}

func main() {
	playerA := makeServer(":3000")
	playerB := makeServer(":4000")
	playerC := makeServer(":5000")

	playerC.Connect(playerA.ListenAddr)

	time.Sleep(2 * time.Millisecond)

	playerB.Connect(playerC.ListenAddr)

	_ = playerA
	_ = playerB

	select {}
}
