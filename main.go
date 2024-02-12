package main

import (
	"fmt"
	"github.com/hmuir28/go-poker/deck"
)

func main() {
	card := deck.NewCard(deck.Spades, 1)
	fmt.Println(card)
}