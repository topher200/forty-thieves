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

// GetGameState returns the latest gamestate for a user
func (db *GameStateDB) GetGameState(userRow UserRow) (*libgame.GameState, error) {
	var gameStateRow GameStateRow
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=$1 ORDER BY id DESC LIMIT 1", db.table)
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

// SaveGameState saves the current gamestate for a user. Does not delete old gamestates.
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
	logrus.Info("Saved new gamestate to db")
	return nil
}

func (db *GameStateDB) DeleteLatestGameState(
	tx *sqlx.Tx, userRow UserRow) error {
	queryWhereStatement := fmt.Sprintf(
		"id=(SELECT id from game_state WHERE user_id=%d ORDER BY id DESC LIMIT 1)",
		userRow.ID)
	res, err := db.DeleteFromTable(tx, queryWhereStatement)
	if err != nil {
		logrus.Warning("Error deleting last game state: ", err)
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected != 1 {
		return errors.New(
			fmt.Sprintf("expected to change 1 row, changed %d", res.RowsAffected))
	}
	logrus.Info("Deleted latest gamestate from db")
	return nil
}
