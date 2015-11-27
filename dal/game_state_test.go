package dal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/topher200/forty-thieves/libgame"
)

func newGameStateDBForTest(t *testing.T) *GameStateDB {
	return NewGameStateDB(newDbForTest(t))
}

func TestGetGameState(t *testing.T) {
	u := newUserForTest(t)
	userRow := u.signupNewUserRowForTest(t)

	gameState := GetGameState(*userRow)
	assert.NotNil(t, gameState)
}

func TestSaveGameState(t *testing.T) {
	u := newUserForTest(t)
	userRow := u.signupNewUserRowForTest(t)
	gameStateDB := newGameStateDBForTest(t)

	gameState := libgame.NewGame()
	err := gameStateDB.SaveGameState(nil, *userRow, gameState)
	assert.Nil(t, err)
}
