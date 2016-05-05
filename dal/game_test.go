package dal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newGameDBForTest(t *testing.T) *GameDB {
	return NewGameDB(newDbForTest(t))
}

func (gameDB *GameDB) createNewGameRowForTest(t *testing.T, userRow UserRow) int64 {
	gameID, err := gameDB.CreateNewGame(nil, userRow)
	assert.Nil(t, err)
	return gameID
}
