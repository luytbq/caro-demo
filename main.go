package main

import (
	"fmt"
	"math/rand"

	"luytbq.com/caro/caro"
)

// create function main
func main() {
	const WIDTH = 5
	const HEIGHT = 5

	game := caro.NewGame(WIDTH, HEIGHT)

	for !game.Over {
		x := rand.Intn(WIDTH - 1)
		y := rand.Intn(HEIGHT - 1)

		fmt.Printf("Player %s move: (%d, %d)\n", game.CurrentPlayer(), x, y)
		game.Move(x, y)
	}

}
