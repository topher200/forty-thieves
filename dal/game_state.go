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
	ID             int64  `db:"id"`
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

func (db *GameStateDB) GetGameState(userRow UserRow) (*libgame.GameState, error) {
	var gameStateRow GameStateRow
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=$1 LIMIT 1", db.table)
	err := db.db.Get(&gameStateRow, query, userRow.ID)
	if err != nil {
		return nil, fmt.Errorf("Error on query: %v", err)
	}
	var gameState libgame.GameState
	decoder := gob.NewDecoder(bytes.NewBuffer(gameStateRow.BinarizedState))
	err = decoder.Decode(&gameState)
	if err != nil {
		return nil, fmt.Errorf("Error decoding: %v", err)
	}
	return &gameState, nil
}

func (db *GameStateDB) SaveGameState(
	tx *sqlx.Tx, userRow UserRow, gameState libgame.GameState) error {
	var binarizedState bytes.Buffer
	encoder := gob.NewEncoder(&binarizedState)
	encoder.Encode(gameState)
	dataStruct := GameStateRow{}
	dataStruct.UserID = userRow.ID
	dataStruct.BinarizedState = binarizedState.Bytes()

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