package deck

import (
	"fmt"
	"strconv"
)

type Suit int

func (s Suit) String() string {
	switch s {
	case Spades: 
		return "SPADES"

	case Harts:
		return "HARTS"

	case Diamonds:
		return "DIAMONDS"
	
	case Clubs:
		return "CLUBS"
	
	default:
		panic("Invalid card suit")

	}
}

const (
	Spades Suit = iota
	Harts
	Diamonds
	Clubs
)

type Card struct {
	suit Suit
	value int
}

func (c Card) String() string {
	value := strconv.Itoa(c.value)

	if c.value == 1 {
		value = "ACE"
	}

	return fmt.Sprintf("%s of %s %s", value, c.suit, suitToUnicode(c.suit))
}

func NewCard(s Suit, v int) Card {
	if v > 13 {
		panic("The value of the card cannot be higher than 13")
	}

	return Card{
		suit: s,
		value: v,
	}
}

type Deck [52]Card

func New() Deck {

	var (
		nSuits = 4
		nCards = 13
		d = [52] Card{}
	)
	
	x := 0
	for i := 0; i < nSuits; i++ {

		for j := 0; j < nCards; j++ {
			d[x] = NewCard(Suit(i), j+1)
			x++
		}

	}

	return d
}

func suitToUnicode(s Suit) string {
	switch s {
		case Spades: 
			return "♠"
	
		case Harts:
			return "♥"
	
		case Diamonds:
			return "♦"
		
		case Clubs:
			return "♣"
		
		default:
			panic("Invalid card suit")
	
	}
}
