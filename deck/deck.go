package deck

type Suit int

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

func NewCard(s Suit, v int) Card {
	if v > 13 {
		panic("The value of the card cannot be higher than 13")
	}

	return Card{
		suit: s,
		value: v,
	}
}
