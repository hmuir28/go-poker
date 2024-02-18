package main

import (
	"log"
	"time"

	"github.com/hmuir28/go-poker/p2p"
)

func main() {
	config := p2p.ServerConfig{
		Version:     "GO POKER v0.1-alpha",
		ListenAddr:  ":3000",
		GameVariant: p2p.Other,
	}
	server := p2p.NewServer(config)

	go server.Start()

	time.Sleep(1 * time.Second)

	remoteConfig := p2p.ServerConfig{
		Version:     "GO POKER v0.1-alpha",
		ListenAddr:  ":4000",
		GameVariant: p2p.TEXAS_HOLDEM,
	}
	remoteServer := p2p.NewServer(remoteConfig)
	go remoteServer.Start()
	if err := remoteServer.Connect(":3000"); err != nil {
		log.Fatal(err)
	}

	select {}
}
