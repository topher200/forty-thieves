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
	ID     int64 `db:"id"`
	UserID int64 `db:"user_id"`
}

func NewGameDB(db *sqlx.DB) *GameDB {
	gs := &GameDB{}
	gs.db = db
	gs.table = "game"
	gs.hasID = true

	return gs
}

// GetLatestGame gets the most recent game (by id) for the given user
//
// Returns error if there are no games for the given user
func (db *GameDB) GetLatestGame(userRow UserRow) (*libgame.Game, error) {
	var gameRow GameRow
	query := fmt.Sprintf(
		"SELECT * FROM %s WHERE user_id=$1 ORDER BY id DESC LIMIT 1", db.table)
	err := db.db.Get(&gameRow, query, userRow.ID)
	if err != nil {
		return nil, fmt.Errorf("Error on query: %v", err)
	}

	var game libgame.Game
	game.ID = gameRow.ID
	return &game, nil
}

// CreateNewGame creates a new game, saves it to the database, and returns it
func (db *GameDB) CreateNewGame(tx *sqlx.Tx, userRow UserRow) (*libgame.Game, error) {
	dataMap := make(map[string]interface{})
	dataMap["user_id"] = userRow.ID
	insertResult, err := db.InsertIntoTable(tx, dataMap)
	if err != nil {
		logrus.Warning("error saving game: ", err)
		return nil, err
	}

	rowsAffected, err := insertResult.RowsAffected()
	if err != nil || rowsAffected != 1 {
		return nil, errors.New(
			fmt.Sprintf("expected to change 1 row, changed %d", insertResult.RowsAffected))
	}

	id, err := insertResult.LastInsertId()
	logrus.Infof("Saved new game (id %d) to db", id)
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
			fmt.Sprintf("expected to change 1 row, changed %d", res.RowsAffected))
	}
	return nil
}
