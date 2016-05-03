package dal

import "github.com/jmoiron/sqlx"

type AnalyzationQueueDB struct {
	Base
}

type AnalyzationQueueRow struct {
	ID          int64 `db:"id"`
	GameStateID int64 `db:"game_state_id"`
	Analyzed    bool  `db:"analyzed"`
}

func NewAnalyzationQueueDB(db *sqlx.DB) *AnalyzationQueueDB {
	gs := &AnalyzationQueueDB{}
	gs.db = db
	gs.table = "analyzation_queue"
	gs.hasID = true

	return gs
}
