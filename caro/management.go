package caro

func NewGame(width, height int) *CaroGame {
	players := [2]Player{
		{sign: SIGN_1},
		{sign: SIGN_2},
	}
	game := &CaroGame{
		width:             width,
		height:            height,
		board:             make([][]string, height),
		players:           players,
		currentPlayer:     players[0],
		winConditionCount: WIN_CONDITION_COUNT,
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
