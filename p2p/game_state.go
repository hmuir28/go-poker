package p2p

type GameStatus uint32

const (
	GameStatusDealing GameStatus = iota
	GameStatusPreflop
	GameStatusFlop
	GameStatusTurn
	GameStatusRiver
)

type GameState struct {
	isDealer   bool
	gameStatus GameStatus
}

func NewGameState() *GameState {
	return &GameState{}
}

func (g *GameState) loop() {

}
