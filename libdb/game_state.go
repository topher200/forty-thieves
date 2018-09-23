package libdb

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"github.com/topher200/forty-thieves/libgame"
)

type GameStateDB struct {
	Base
}

type GameStateRow struct {
	GameStateID       uuid.UUID     `db:"game_state_id"`
	PreviousGameState uuid.NullUUID `db:"previous_game_state"`
	GameID            int64         `db:"game_id"`
	MoveNum           int64         `db:"move_num"`
	Score             int           `db:"score"`
	Status            string        `db:"status"`
	BinarizedState    []byte        `db:"binarized_state"`
}

func NewGameStateDB(db *sqlx.DB) *GameStateDB {
	gs := &GameStateDB{}
	gs.db = db
	gs.table = "game_state"
	gs.hasID = false

	return gs
}

// GetGameStateById returns the game state for the given id
//
// Returns error if there are no game states for the given game
func (db *GameStateDB) GetGameStateById(gameStateID uuid.UUID) (*libgame.GameState, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE game_state_id=$1 LIMIT 1", db.table)
	gameState, err := db.getSingleGameState(query, gameStateID.String())
	if err != nil {
		return nil, fmt.Errorf("Error getting gamestate by id: %v", err)
	}
	return gameState, nil
}

// GetFirstGameState returns the first gamestate for the given game
//
// Returns error if there are no game states for the given game
func (db *GameStateDB) GetFirstGameState(game libgame.Game) (*libgame.GameState, error) {
	query := fmt.Sprintf(
		"SELECT * FROM %s WHERE game_id=$1 and move_num=0 LIMIT 1", db.table)
	gameState, err := db.getSingleGameState(query, strconv.FormatInt(game.ID, 10))
	if err != nil {
		return nil, fmt.Errorf("Error getting first gamestate: %v", err)
	}
	return gameState, nil
}

// GetNextToAnalyze returns the highest priority GameState to analyze.
//
// Returns the unprocessed GameState from the given game with the lowest score
// (primary sort) and the fewest number of moves (secondary sort).
func (db *GameStateDB) GetNextToAnalyze(game libgame.Game) (*libgame.GameState, error) {
	query := fmt.Sprintf(`
	    UPDATE game_state SET status='CLAIMED'
	    WHERE game_state_id = (
		SELECT game_state_id FROM game_state
		WHERE game_id=$1 AND status='UNPROCESSED'
		ORDER BY score ASC, move_num ASC
		LIMIT 1
		FOR UPDATE SKIP LOCKED
	    )
	    RETURNING *
	`)
	var gameStateRow GameStateRow
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

// getSingleGameState is a helper function for getting and parsing a game state
//
// Implementation note: there's no reason why this function can't take more than
// one query argument. I just implemented the first rev taking a single one for
// convenience.
func (db *GameStateDB) getSingleGameState(query string, arg string) (*libgame.GameState, error) {
	var gameStateRow GameStateRow
	err := db.db.Get(&gameStateRow, query, arg)
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

// GetChildGameStates queries for the UUIDs of all the game states that are children of the given one
func (db *GameStateDB) GetChildGameStates(gameState libgame.GameState) ([]uuid.UUID, error) {
	query := fmt.Sprintf("SELECT game_state_id FROM %s WHERE game_id=$1 and previous_game_state=$2", db.table)
	var childIds []uuid.UUID
	err := db.db.Select(&childIds, query, gameState.GameID, gameState.GameStateID)
	if err != nil {
		return nil, fmt.Errorf("Error on query: %v", err)
	}

	return childIds, nil
}

// SaveGameState saves the given gamestate to the db given the game and the gamestate
func (db *GameStateDB) SaveGameState(tx *sqlx.Tx, gameState libgame.GameState) error {
	var binarizedState bytes.Buffer
	encoder := gob.NewEncoder(&binarizedState)
	encoder.Encode(gameState)
	dataStruct := GameStateRow{}
	dataStruct.GameID = gameState.GameID
	dataStruct.BinarizedState = binarizedState.Bytes()
	dataStruct.GameStateID = gameState.GameStateID
	dataStruct.MoveNum = gameState.MoveNum
	dataStruct.PreviousGameState = gameState.PreviousGameState
	dataStruct.Score = gameState.Score

	dataMap := make(map[string]interface{})
	dataMap["game_id"] = dataStruct.GameID
	dataMap["binarized_state"] = dataStruct.BinarizedState
	dataMap["game_state_id"] = dataStruct.GameStateID
	dataMap["previous_game_state"] = dataStruct.PreviousGameState
	dataMap["move_num"] = dataStruct.MoveNum
	dataMap["score"] = dataStruct.Score
	insertResult, err := db.InsertIntoTable(tx, dataMap)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"gamesState": gameState,
			"err":        err,
		}).Warning("error saving game state")
		return err
	}
	rowsAffected, err := insertResult.RowsAffected()
	if err != nil || rowsAffected != 1 {
		return errors.New(
			fmt.Sprintf("expected to change 1 row, changed %d", rowsAffected))
	}

	logrus.WithFields(logrus.Fields{
		"id": gameState.GameStateID,
	}).Info("saved new gamestate to db")
	return nil
}

func (db *GameStateDB) MarkAsProcessed(tx *sqlx.Tx, gameState libgame.GameState) error {
	res, err := db.db.Exec(
		"UPDATE game_state SET status='PROCESSED' WHERE game_state_id=$1",
		gameState.GameStateID)
	if err != nil {
		logrus.Warning("Error updating game state: ", err)
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected != 1 {
		return errors.New(
			fmt.Sprintf("expected to change 1 row, changed %d", rowsAffected))
	}
	return nil
}

// DeleteGameState deletes the given gamestate
func (db *GameStateDB) DeleteGameState(
	tx *sqlx.Tx, gameState libgame.GameState) error {
	queryWhereStatement := fmt.Sprintf("game_state_id='%v'", gameState.GameStateID)
	res, err := db.DeleteFromTable(tx, queryWhereStatement)
	if err != nil {
		logrus.Warning("Error deleting game state: ", err)
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected != 1 {
		return errors.New(
			fmt.Sprintf("expected to change 1 row, changed %d", rowsAffected))
	}
	return nil
}
