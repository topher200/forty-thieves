package dal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/topher200/forty-thieves/libgame"
)

func newGameStateDBForTest(t *testing.T) *GameStateDB {
	return NewGameStateDB(newDbForTest(t))
}

func TestGetEmptyGameState(t *testing.T) {
	u := NewUserDBForTest(t)
	userRow := u.signupNewUserRowForTest(t)
	gameDB := newGameDBForTest(t)
	gameID := gameDB.createNewGameRowForTest(t, *userRow)
	gameStateDB := newGameStateDBForTest(t)

	// We should err, since we haven't set a game yet
	_, err := gameStateDB.GetLatestGameState(gameID)
	assert.NotNil(t, err)
}

func TestSaveAndGetGameState(t *testing.T) {
	u := NewUserDBForTest(t)
	userRow := u.signupNewUserRowForTest(t)
	gameDB := newGameDBForTest(t)
	gameID := gameDB.createNewGameRowForTest(t, *userRow)
	gameStateDB := newGameStateDBForTest(t)

	// Save game state
	originalGameState := libgame.NewGame()
	err := gameStateDB.SaveGameState(nil, gameID, 0, originalGameState, 0)
	assert.Nil(t, err)

	// Retrieved saved game state
	retrievedGameState, err := gameStateDB.GetLatestGameState(gameID)
	assert.Nil(t, err)
	assert.Equal(t, originalGameState, *retrievedGameState)

	// Now delete that game state
	err = gameStateDB.DeleteLatestGameState(nil, gameID)
	assert.Nil(t, err)
}
