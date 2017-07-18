package dal

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"github.com/topher200/forty-thieves/libgame"
)

type GameStateDB struct {
	Base
}

type GameStateRow struct {
	ID                int64     `db:"id"`
	GameStateID       uuid.UUID `db:"game_state_id"`
	PreviousGameState uuid.UUID `db:"previous_game_state"`
	GameID            int64     `db:"game_id"`
	MoveNum           int64     `db:"move_num"`
	Score             int       `db:"score"`
	BinarizedState    []byte    `db:"binarized_state"`
}

func NewGameStateDB(db *sqlx.DB) *GameStateDB {
	gs := &GameStateDB{}
	gs.db = db
	gs.table = "game_state"
	gs.hasID = true

	return gs
}

// GetGameState returns the latest gamestate for the given game
//
// Returns error if there are no game states for the given game
func (db *GameStateDB) GetGameState(game libgame.Game) (*libgame.GameState, error) {
	var gameStateRow GameStateRow
	query := fmt.Sprintf(
		"SELECT * FROM %s WHERE game_id=$1 ORDER BY id DESC LIMIT 1", db.table)
	err := db.db.Get(&gameStateRow, query, game.ID)
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

// SaveGameState saves the given gamestate to the db given the game and the gamestate
func (db *GameStateDB) SaveGameState(
	tx *sqlx.Tx, game libgame.Game, gameState libgame.GameState) error {
	var binarizedState bytes.Buffer
	encoder := gob.NewEncoder(&binarizedState)
	encoder.Encode(gameState)
	dataStruct := GameStateRow{}
	dataStruct.GameID = game.ID
	dataStruct.BinarizedState = binarizedState.Bytes()
	dataStruct.GameStateID = gameState.GameStateID
	dataStruct.PreviousGameState = gameState.PreviousGameState
	dataStruct.Score = gameState.Score

	dataMap := make(map[string]interface{})
	dataMap["game_id"] = dataStruct.GameID
	dataMap["binarized_state"] = dataStruct.BinarizedState
	dataMap["game_state_id"] = dataStruct.GameStateID
	dataMap["previous_game_state"] = dataStruct.PreviousGameState
	dataMap["score"] = dataStruct.Score
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

// DeleteGameState deletes the given gamestate
func (db *GameStateDB) DeleteGameState(
	tx *sqlx.Tx, gameState libgame.GameState) error {
	queryWhereStatement := fmt.Sprintf("id=%d", gameState.ID)
	res, err := db.DeleteFromTable(tx, queryWhereStatement)
	if err != nil {
		logrus.Warning("Error deleting game state: ", err)
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected != 1 {
		return errors.New(
			fmt.Sprintf("expected to change 1 row, changed %d", res.RowsAffected))
	}
	return nil
}
