package main

import (
	"github.com/hmuir28/go-poker/p2p"
)

func main() {
	// rand.Seed(time.Now().UnixNano())

	// for i := 0; i < 10; i++ {
	// 	d := deck.New()
	// 	fmt.Println(d)
	// 	fmt.Println()
	// }

	config := p2p.ServerConfig{
		ListenAddr: ":3000",
	}
	server := p2p.NewServer(config)

	server.Start()

}
