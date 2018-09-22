package libdb

import (
	"errors"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	"github.com/topher200/forty-thieves/libgame"
)

type GameDB struct {
	Base
}

type GameRow struct {
	ID int64 `db:"id"`
}

func NewGameDB(db *sqlx.DB) *GameDB {
	gs := &GameDB{}
	gs.db = db
	gs.table = "game"
	gs.hasID = true

	return gs
}

// GetLatestGame gets the most recent game (by id)
//
// Returns error if there are no games
func (db *GameDB) GetLatestGame() (*libgame.Game, error) {
	var gameRow GameRow
	query := fmt.Sprintf(
		"SELECT * FROM %s ORDER BY id DESC LIMIT 1", db.table)
	err := db.db.Get(&gameRow, query)
	if err != nil {
		return nil, fmt.Errorf("Error on query: %v", err)
	}

	var game libgame.Game
	game.ID = gameRow.ID
	return &game, nil
}

// CreateNewGame creates a new game, saves it to the database, and returns it
func (db *GameDB) CreateNewGame(tx *sqlx.Tx) (*libgame.Game, error) {
	dataMap := make(map[string]interface{})
	insertResult, err := db.InsertIntoTable(tx, dataMap)
	if err != nil {
		logrus.Warning("error saving game: ", err)
		return nil, err
	}

	rowsAffected, err := insertResult.RowsAffected()
	if err != nil || rowsAffected != 1 {
		return nil, errors.New(
			fmt.Sprintf("expected to change 1 row, changed %d", rowsAffected))
	}

	id, err := insertResult.LastInsertId()
	logrus.WithFields(logrus.Fields{
		"id": id,
	}).Info("saved new game to db")
	var game libgame.Game
	game.ID = id
	return &game, nil
}

// DeleteGame deletes the given libgame.Game
func (db *GameDB) DeleteGame(tx *sqlx.Tx, game libgame.Game) error {
	queryWhereStatement := fmt.Sprintf("id=%d", game.ID)
	res, err := db.DeleteFromTable(tx, queryWhereStatement)
	if err != nil {
		logrus.Warning("Error deleting game: ", err)
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected != 1 {
		return errors.New(
			fmt.Sprintf("expected to change 1 row, changed %d", rowsAffected))
	}
	return nil
}
