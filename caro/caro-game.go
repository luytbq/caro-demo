package caro

import (
	"bytes"
	"fmt"
	"time"
)

type CaroGame struct {
	width             int
	height            int
	board             [][]string
	players           [2]Player
	currentPlayer     Player
	winConditionCount int
	Over              bool
	HasWinner         bool
	moves             int
	totalTiles        int
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

func (g *CaroGame) PrintBoard() {
	buffer := bytes.Buffer{}

	// Warning: no warranty for other OS
	// clearScreen := func() {
	// 	fmt.Print("\033[H\033[2J")
	// }
	// clearScreen()

	addHLine := func() {
		for j := 0; j < g.width; j++ {
			buffer.WriteString(H_DASH)
			if j == g.width-1 {
				buffer.WriteString(NEW_LINE)
			}
		}
	}

	// begin printing
	buffer.WriteString(NEW_LINE)
	for i := 0; i < g.height; i++ {
		addHLine()
		buffer.WriteString(V_DASH)
		for j := 0; j < g.width; j++ {
			buffer.WriteString(fmt.Sprintf(" %s %s", g.board[i][j], V_DASH))
			if j == g.width-1 {
				buffer.WriteString(NEW_LINE)
			}
		}
	}
	addHLine()

	fmt.Print(buffer.String())
}

func (g *CaroGame) Move(x int, y int) {
	if err := g.validateMove(x, y); err != nil {
		fmt.Print(err.Error())
		return
	}

	g.board[y][x] = string(g.currentPlayer.sign)
	g.moves++
	g.PrintBoard()

	g.checkWin(x, y)
	if g.HasWinner {
		return
	}

	g.nextTurn()
	time.Sleep(700 * time.Millisecond)
}
func (g *CaroGame) nextTurn() {
	if g.currentPlayer == g.players[0] {
		g.currentPlayer = g.players[1]
	} else {
		g.currentPlayer = g.players[0]
	}
}

func (g *CaroGame) validateMove(x int, y int) error {
	if !g.validTile(x, y) {
		return fmt.Errorf("invalid move (%d, %d): out of range", x, y)
	}
	if g.board[y][x] != EMPTY {
		return fmt.Errorf("invalid move (%d, %d): not an empty tile", x, y)
	}
	return nil
}

func (g *CaroGame) validTile(x int, y int) bool {
	return x >= 0 && x < g.width && y >= 0 && y < g.height
}

func (g *CaroGame) checkWin(x int, y int) {
	// c := make(chan bool)
	// wg := sync.WaitGroup{}

	// f := func(vectorX int, vectorY int) {
	// 	defer wg.Done()
	// 	c <- g.checkWinByVector(x, y, vectorX, vectorY)
	// }

	// wg.Add(4)
	// go f(1, 0)
	// go f(0, 1)
	// go f(1, 1)
	// go f(1, -1)

	// wg.Wait()
	// close(c)

	// for hasWinner := range c {
	// 	g.HasWinner = hasWinner
	// 	break
	// }

	g.HasWinner = g.checkWinByVector(x, y, 1, 0) || g.checkWinByVector(x, y, 0, 1) || g.checkWinByVector(x, y, 1, 1) || g.checkWinByVector(x, y, 1, -1)
	if g.HasWinner || g.moves == g.totalTiles {
		g.Over = true
	}

	// debug
	// if g.Over {
	// 	fmt.Printf("checkWinByVector(x, y, 1, 0) %t\n", g.checkWinByVector(x, y, 1, 0))
	// 	fmt.Printf("checkWinByVector(x, y, 0, 1) %t\n", g.checkWinByVector(x, y, 0, 1))
	// 	fmt.Printf("checkWinByVector(x, y, 1, 1) %t\n", g.checkWinByVector(x, y, 1, 1))
	// 	fmt.Printf("checkWinByVector(x, y, -1, -1) %t\n", g.checkWinByVector(x, y, -1, -1))
	// }

}

func (g *CaroGame) checkWinByVector(x int, y int, vectorX int, vectorY int) bool {
	maxForward := 0
	maxBackword := 0

	for maxForward <= g.winConditionCount {
		checkX := x + maxForward*vectorX
		checkY := y + maxForward*vectorY
		if g.validTile(checkX, checkY) && g.board[checkY][checkX] == g.currentPlayer.sign {
			maxForward++
		} else {
			break
		}
	}

	for maxBackword <= g.winConditionCount {
		checkX := x - maxBackword*vectorX
		checkY := y - maxBackword*vectorY
		if g.validTile(checkX, checkY) && g.board[checkY][checkX] == g.currentPlayer.sign {
			maxBackword++
		} else {
			break
		}
	}
	hasWinner := maxBackword >= g.winConditionCount || maxForward >= g.winConditionCount
	if hasWinner {
		fmt.Printf("vectorX: %d, vectorY: %d, maxForward: %d, maxBackword: %d\n", vectorX, vectorY, maxForward, maxBackword)
	}
	return hasWinner
}

func (g *CaroGame) CurrentPlayer() string {
	return g.currentPlayer.sign
}
