package dal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/topher200/forty-thieves/libgame"
)

func newGameStateDBForTest(t *testing.T) *GameStateDB {
	return NewGameStateDB(newDbForTest(t))
}

func setupNewGameForTest(t *testing.T, gameDB GameDB) *libgame.Game {
	u := NewUserDBForTest(t)
	userRow := u.signupNewUserRowForTest(t)

	game, err := gameDB.CreateNewGame(nil, *userRow)
	assert.Nil(t, err)
	return game
}

func TestGetEmptyGameState(t *testing.T) {
	gameStateDB := newGameStateDBForTest(t)
	gameDB := newGameDBForTest(t)
	game := setupNewGameForTest(t, *gameDB)
	defer gameDB.DeleteGame(nil, *game)

	// We should err, since we haven't set a game yet
	_, err := gameStateDB.GetLatestGameState(*game)
	assert.NotNil(t, err)
}

func TestSaveAndGetGameState(t *testing.T) {
	gameStateDB := newGameStateDBForTest(t)
	gameDB := newGameDBForTest(t)
	game := setupNewGameForTest(t, *gameDB)
	defer gameDB.DeleteGame(nil, *game)

	// Save game state
	originalGameState := libgame.DealNewGame(*game)
	err := gameStateDB.SaveGameState(nil, originalGameState)
	defer gameStateDB.DeleteGameState(nil, originalGameState)
	assert.Nil(t, err)

	// Retrieved saved game state by game
	retrievedGameState, err := gameStateDB.GetLatestGameState(*game)
	assert.Nil(t, err)
	assert.Equal(t, originalGameState, *retrievedGameState)

	// Retrieved saved game state by id
	retrievedGameState, err = gameStateDB.GetGameStateById(originalGameState.GameStateID)
	assert.Nil(t, err)
	assert.Equal(t, originalGameState, *retrievedGameState)
}
