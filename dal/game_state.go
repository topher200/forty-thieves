package dal

import (
	"github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	"github.com/topher200/forty-thieves/libgame"
)

type GameStateDB struct {
	Base
}

type GameStateRow struct {
	UserID          int64  `db:"user_id"`
	SerializedState string `db:"serialized_state"`
}

func NewGameStateDB(db *sqlx.DB) *GameStateDB {
	gs := &GameStateDB{}
	gs.db = db
	gs.table = "game_state"
	gs.hasID = true

	return gs
}

func GetGameState(userRow UserRow) *libgame.GameState {
	// TODO(topher)
	return nil
}

func (db *GameStateDB) SaveGameState(
	tx *sqlx.Tx, userRow UserRow, gameState libgame.GameState) error {
	data := make(map[string]interface{})
	data["user_id"] = userRow.ID
	// TODO(topher)
	data["serialized_state"] = "asdf"

	_, err := db.InsertIntoTable(tx, data)
	if err == nil {
		logrus.Warning("error saving game state:", err)
		return err
	} else {
		logrus.Info("Success saving game state!")
		return nil
	}
}
