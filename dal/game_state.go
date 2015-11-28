package dal

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	"github.com/topher200/forty-thieves/libgame"
)

type GameStateDB struct {
	Base
}

type GameStateRow struct {
	UserID         int64  `db:"user_id"`
	BinarizedState []byte `db:"binarized_state"`
}

func NewGameStateDB(db *sqlx.DB) *GameStateDB {
	gs := &GameStateDB{}
	gs.db = db
	gs.table = "game_state"
	gs.hasID = true

	return gs
}

func GetGameState(userRow UserRow) *libgame.GameState {
	// TODO(topher)
	return nil
}

func (db *GameStateDB) SaveGameState(
	tx *sqlx.Tx, userRow UserRow, gameState libgame.GameState) error {
	var binarizedState bytes.Buffer
	encoder := gob.NewEncoder(&binarizedState)
	encoder.Encode(gameState)
	dataStruct := GameStateRow{userRow.ID, binarizedState.Bytes()}

	dataMap := make(map[string]interface{})
	dataMap["user_id"] = dataStruct.UserID
	dataMap["binarized_state"] = dataStruct.BinarizedState
	insertResult, err := db.InsertIntoTable(tx, dataMap)
	if err != nil {
		logrus.Warning("error saving game state:", err)
		return err
	}
	rowsAffected, err := insertResult.RowsAffected()
	if err != nil || rowsAffected != 1 {
		return errors.New(
			fmt.Sprintf("expected to change 1 row, changed %d", insertResult.RowsAffected))
	}
	return nil
}
