package dal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newGameDBForTest(t *testing.T) *GameDB {
	return NewGameDB(newDbForTest(t))
}

func TestGetEmptyGame(t *testing.T) {
	u := NewUserDBForTest(t)
	userRow := u.signupNewUserRowForTest(t)
	gameDB := newGameDBForTest(t)

	// We should err, since we haven't set a game yet
	_, err := gameDB.GetLatestGame(*userRow)
	assert.NotNil(t, err)
}

func TestCreateAndGetGame(t *testing.T) {
	u := NewUserDBForTest(t)
	userRow := u.signupNewUserRowForTest(t)
	gameDB := newGameDBForTest(t)

	// Create new game
	originalGame, err := gameDB.CreateNewGame(nil, *userRow)
	assert.Nil(t, err)

	// Now retrieve our new game
	retrievedGame, err := gameDB.GetLatestGame(*userRow)
	assert.Nil(t, err)
	assert.Equal(t, *originalGame, *retrievedGame)
}
