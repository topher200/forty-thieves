package dal

import (
	"errors"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
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

func (db *GameDB) GetLatestGameID(userRow UserRow) (int64, error) {
	var gameRow GameRow
	query := fmt.Sprintf(
		"SELECT * FROM %s WHERE user_id=$1 ORDER BY id DESC LIMIT 1", db.table)
	err := db.db.Get(&gameRow, query, userRow.ID)
	if err != nil {
		return 0, fmt.Errorf("Error on query: %v", err)
	}
	return gameRow.ID, nil
}

// CreateNewGame adds a new game to the database. Returns the game's ID
func (db *GameDB) CreateNewGame(tx *sqlx.Tx, userRow UserRow) (int64, error) {
	dataMap := make(map[string]interface{})
	dataMap["user_id"] = userRow.ID
	insertResult, err := db.InsertIntoTable(tx, dataMap)
	if err != nil {
		logrus.Warning("error saving game: ", err)
		return 0, err
	}

	rowsAffected, err := insertResult.RowsAffected()
	if err != nil || rowsAffected != 1 {
		return 0, errors.New(
			fmt.Sprintf("expected to change 1 row, changed %d", insertResult.RowsAffected))
	}

	id, err := insertResult.LastInsertId()
	logrus.Infof("Saved new game (id %d) to db", id)
	return id, nil
}
