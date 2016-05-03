package dal

import (
	"fmt"

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
