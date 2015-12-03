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
	u := NewUserForTest(t)
	userRow := u.signupNewUserRowForTest(t)
	gameStateDB := newGameStateDBForTest(t)

	// We should err, since we haven't set a game yet
	_, err := gameStateDB.GetGameState(*userRow)
	assert.NotNil(t, err)
}

func TestSaveAndGetGameState(t *testing.T) {
	u := NewUserForTest(t)
	userRow := u.signupNewUserRowForTest(t)
	gameStateDB := newGameStateDBForTest(t)

	// Save game state
	originalGameState := libgame.NewGame()
	err := gameStateDB.SaveGameState(nil, *userRow, originalGameState)
	assert.Nil(t, err)

	// Retrieved saved game state
	retrievedGameState, err := gameStateDB.GetGameState(*userRow)
	assert.Nil(t, err)
	assert.Equal(t, originalGameState, *retrievedGameState)
}
