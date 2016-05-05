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
	ID                  int64  `db:"id"`
	GameID              int64  `db:"game_id"`
	MoveNum             int64  `db:"move_num"`
	Score               int64  `db:"score"`
	BinarizedState      []byte `db:"binarized_state"`
	PreviousGameStateID int64  `db:"previous_game_state_id"`
}

func NewGameStateDB(db *sqlx.DB) *GameStateDB {
	gs := &GameStateDB{}
	gs.db = db
	gs.table = "game_state"
	gs.hasID = true

	return gs
}

// GetLatestGameState returns the latest gamestate for a game
func (db *GameStateDB) GetLatestGameState(gameID int64) (*libgame.GameState, error) {
	var gameStateRow GameStateRow
	query := fmt.Sprintf(
		"SELECT * FROM %s WHERE game_id=$1 ORDER BY id DESC LIMIT 1", db.table)
	err := db.db.Get(&gameStateRow, query, gameID)
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

// SaveGameState saves the current gamestate for a game. Does not delete old gamestates.
func (db *GameStateDB) SaveGameState(
	tx *sqlx.Tx, gameID int64, moveNum int64, gameState libgame.GameState,
	previousGameStateID int64) error {
	var binarizedState bytes.Buffer
	encoder := gob.NewEncoder(&binarizedState)
	encoder.Encode(gameState)
	dataStruct := GameStateRow{}
	dataStruct.GameID = gameID
	dataStruct.MoveNum = moveNum
	dataStruct.Score = int64(gameState.Score)
	dataStruct.BinarizedState = binarizedState.Bytes()
	dataStruct.PreviousGameStateID = previousGameStateID

	dataMap := make(map[string]interface{})
	dataMap["game_id"] = dataStruct.GameID
	dataMap["move_num"] = dataStruct.MoveNum
	dataMap["score"] = dataStruct.Score
	dataMap["binarized_state"] = dataStruct.BinarizedState
	dataMap["previous_game_state_id"] = dataStruct.PreviousGameStateID
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

	id, err := insertResult.LastInsertId()
	logrus.Infof("Saved new gamestate (id %d) to db", id)
	return nil
}

func (db *GameStateDB) DeleteLatestGameState(
	tx *sqlx.Tx, gameID int64) error {
	// Get latest gamestate's ID
	var gameStateRow GameStateRow
	query := fmt.Sprintf(
		"SELECT * FROM %s WHERE game_id=$1 ORDER BY id DESC LIMIT 1", db.table)
	err := db.db.Get(&gameStateRow, query, gameID)
	if err != nil {
		return fmt.Errorf("Error getting latest gamestate: %v", err)
	}
	logrus.Infof("Deleting latest gamestate (id %d) from db", gameStateRow.ID)

	// Delete the gamestate
	queryWhereStatement := fmt.Sprintf("id=%d", gameStateRow.ID)
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
	return nil
}
