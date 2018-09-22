package libdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/topher200/forty-thieves/libgame"
)

func newGameDBForTest(t *testing.T) *GameDB {
	return NewGameDB(newDbForTest(t))
}

func CreateNewGameForTest(t *testing.T) *libgame.Game {
	gameDB := newGameDBForTest(t)

	game, err := gameDB.CreateNewGame(nil)
	assert.Nil(t, err)
	return game
}

func TestCreateAndDeleteNewGame(t *testing.T) {
	// Create a new game
	game := CreateNewGameForTest(t)

	// Delete it
	gameDB := newGameDBForTest(t)
	err := gameDB.DeleteGame(nil, *game)
	assert.Nil(t, err)
}

func TestGetNewlyCreatedGame(t *testing.T) {
	// Create a new game
	gameDB := newGameDBForTest(t)
	originalGame, err := gameDB.CreateNewGame(nil)
	defer gameDB.DeleteGame(nil, *originalGame)
	assert.Nil(t, err)

	// Now retrieve the game and compare
	retrievedGame, err := gameDB.GetLatestGame()
	assert.Nil(t, err)
	assert.Equal(t, *originalGame, *retrievedGame)
}
