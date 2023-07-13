package main

import (
	"fmt"
)

func main() {
	game := Game{}

	game.AddPlayer(Player{ID: "1", Name: "Player 1"})
	game.AddPlayer(Player{ID: "2", Name: "Player 2"})

	fmt.Println("It's", game.CurrentPlayer().Name, "'s turn")

	game.NextTurn()
	fmt.Println("It's", game.CurrentPlayer().Name, "'s turn")
}
