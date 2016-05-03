package dal

import "github.com/jmoiron/sqlx"

type WinningMovesDB struct {
	Base
}

type WinningMovesRow struct {
	ID          int64 `db:"id"`
	GameStateID int64 `db:"game_state_id"`
}

func NewWinningMovesDB(db *sqlx.DB) *WinningMovesDB {
	gs := &WinningMovesDB{}
	gs.db = db
	gs.table = "winning_moves"
	gs.hasID = true

	return gs
}
