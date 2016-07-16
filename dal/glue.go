package dal

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/topher200/forty-thieves/libgame"
)

func SaveGameState(
	tx *sqlx.Tx, gameDB *GameDB, gameStateDB *GameStateDB, user UserRow,
	newGameState libgame.GameState) error {
	// Get the game
	gameID, err := gameDB.GetLatestGameID(user)
	if err != nil {
		return fmt.Errorf("Error getting game: %v", err)
	}

	// Get the previous game state
	currentGameState, err := gameStateDB.GetLatestGameState(gameID)
	if err != nil {
		return fmt.Errorf("Error getting latest gamestate: %v", err)
	}

	// Save the new game state
	err = gameStateDB.SaveGameState(nil, gameID, currentGameState.MoveNum+1, newGameState,
		currentGameState.ID)
	if err != nil {
		return fmt.Errorf("Error saving gamestate: %v", err)
	}

	return nil
}
