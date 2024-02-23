package main

import (
	"fmt"
	"time"
)

type Game struct {
	width         int
	height        int
	board         [][]string
	players       [2]Player
	currentPlayer Player
}

type Tile struct {
	x, y int
}

type SearchVector struct {
	x, y int
}

type Player struct {
	sign string
}

const (
	EMPTY  = " "
	SIGN_1 = "X"
	SIGN_2 = "O"
)

func NewGame(width, height int) *Game {
	players := [2]Player{
		{sign: SIGN_1},
		{sign: SIGN_2},
	}
	game := &Game{
		width:         width,
		height:        height,
		board:         make([][]string, height),
		players:       players,
		currentPlayer: players[0],
	}
	for i := range game.board {
		game.board[i] = make([]string, width)
	}
	for i := range game.board {
		for j := range game.board[i] {
			game.board[i][j] = " "
		}
	}
	return game
}

func (g *Game) PrintBoard() {
	const V_DASH = "|"
	const H_DASH = "----"

	// Warning: no warranty for other OS
	clearScreen := func() {
		fmt.Print("\033[H\033[2J")
	}
	clearScreen()

	printHLine := func() {
		for j := 0; j < g.width; j++ {
			fmt.Print(H_DASH)
			if j == g.width-1 {
				fmt.Println()
			}
		}
	}

	// begin printing
	fmt.Println()
	for i := 0; i < g.height; i++ {
		printHLine()
		fmt.Print(V_DASH)
		for j := 0; j < g.width; j++ {
			fmt.Printf(" %s %s", g.board[i][j], V_DASH)
			if j == g.width-1 {
				fmt.Printf("\n")
			}
		}
	}
	printHLine()
}

func (g *Game) move(x int, y int) {
	if err := g.validMove(x, y); err != nil {
		fmt.Printf("Invalid move (%d, %d)\n", x, y)
	} else {
		g.board[y][x] = string(g.currentPlayer.sign)
		g.switchPlayer()
		g.PrintBoard()
		time.Sleep(1 * time.Second)
	}
}
func (g *Game) switchPlayer() {
	if g.currentPlayer == g.players[0] {
		g.currentPlayer = g.players[1]
	} else {
		g.currentPlayer = g.players[0]
	}
}

func (g *Game) validMove(x int, y int) error {
	if x < 0 || x >= g.width || y < 0 || y >= g.height {
		return fmt.Errorf("Invalid move (%d, %d): out of range", x, y)
	}
	if g.board[x][y] != EMPTY {
		return fmt.Errorf("Invalid move (%d, %d): not an empty tile", x, y)
	}
	return nil
}

// create function main
func main() {
	const WIDTH = 8
	const HEIGHT = 8

	game := NewGame(WIDTH, HEIGHT)

	game.move(3, 5)
	game.move(2, 4)
	game.move(4, 5)
	game.move(3, 4)
	game.move(4, 5)
	game.move(4, 4)
	game.move(5, 5)
	game.move(4, 6)
	game.move(6, 5)
}
