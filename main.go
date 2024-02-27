package main

import (
	"fmt"

	"luytbq.com/caro/caro"
)

// create function main
func main() {
	const WIDTH = 3
	const HEIGHT = 3

	game := caro.NewGame(WIDTH, HEIGHT)

	moves := make(map[int][2]int)
	moves[0] = [2]int{0, 0}
	moves[1] = [2]int{0, 1}
	moves[2] = [2]int{0, 2}
	moves[3] = [2]int{1, 1}
	moves[4] = [2]int{1, 0}
	moves[5] = [2]int{1, 2}
	moves[6] = [2]int{2, 1}
	moves[7] = [2]int{2, 2}
	moves[8] = [2]int{2, 0}

	fmt.Print(len(moves))

	i := 0
	for !game.Over {
		x := moves[i][0]
		y := moves[i][1]
		delete(moves, i)
		i++

		fmt.Printf("Player %s move: (%d, %d)\n", game.CurrentPlayer(), x, y)
		game.Move(x, y)

		if len(moves) == 0 {
			break
		}

	}
	if !game.HasWinner {
		fmt.Println("It's a draw!")
	} else {
		fmt.Printf("Player %s wins!\n", game.CurrentPlayer())
	}

}
