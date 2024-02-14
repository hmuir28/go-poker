package main

import (
	"time"

	"github.com/hmuir28/go-poker/p2p"
)

func main() {
	config := p2p.ServerConfig{
		Version:    "GO POKER v0.1-alpha",
		ListenAddr: ":3000",
	}
	server := p2p.NewServer(config)

	go server.Start()

	time.Sleep(1 * time.Second)

	remoteConfig := p2p.ServerConfig{
		Version:    "GO POKER v0.1-alpha",
		ListenAddr: ":4000",
	}
	remoteServer := p2p.NewServer(remoteConfig)
	go remoteServer.Start()
	remoteServer.Connect(":3000")

	select {}
}
