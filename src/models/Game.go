package models

type Game struct {
	ID                 string   `json:"id"`
	Players            []Player `json:"players"`
	CurrentPlayerIndex int      `json:"currentPlayerIndex"`
}

func (g *Game) AddPlayer(player Player) {
	g.Players = append(g.Players, player)
}

func (g *Game) NextTurn() {
	g.CurrentPlayerIndex = (g.CurrentPlayerIndex + 1) % len(g.Players)
}

func (g *Game) CurrentPlayer() Player {
	return g.Players[g.CurrentPlayerIndex]
}
