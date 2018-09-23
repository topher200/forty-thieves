package libdb

import (
	"encoding/json"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/topher200/deck"
	"github.com/topher200/forty-thieves/libgame"
)

type decksJSONStruct struct {
	Stock       deck.Deck
	Foundations []deck.Deck
	Tableaus    []deck.Deck
	Waste       deck.Deck
}

// UnmarshalGameState unmarshalls a GameStateRow into a GameState.
func UnmarshalGameState(gameStateRow GameStateRow) (*libgame.GameState, error) {
	var gameState libgame.GameState
	gameState.GameID = gameStateRow.GameID
	gameState.GameStateID = gameStateRow.GameStateID
	gameState.PreviousGameState = gameStateRow.PreviousGameState
	gameState.MoveNum = gameStateRow.MoveNum
	gameState.Score = gameStateRow.Score

	// unmarshal json
	var deckData decksJSONStruct
	err := gameStateRow.DecksJSON.Unmarshal(&deckData)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling gameStateRow: %v", err)
	}
	gameState.Stock = deckData.Stock
	gameState.Foundations = deckData.Foundations
	gameState.Tableaus = deckData.Tableaus
	gameState.Waste = deckData.Waste

	return &gameState, nil
}

func MarshalGameState(gameState libgame.GameState) (*GameStateRow, error) {
	// convert decks to JSON
	decksJSON := decksJSONStruct{
		gameState.Stock,
		gameState.Foundations,
		gameState.Tableaus,
		gameState.Waste,
	}
	decksJSONSerialized, err := json.Marshal(&decksJSON)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"gamesState": gameState,
			"err":        err,
		}).Error("error JSON-ing deck")
		return nil, err
	}

	var gameStateRow GameStateRow
	gameStateRow.GameID = gameState.GameID
	gameStateRow.GameStateID = gameState.GameStateID
	gameStateRow.MoveNum = gameState.MoveNum
	gameStateRow.PreviousGameState = gameState.PreviousGameState
	gameStateRow.Score = gameState.Score
	gameStateRow.DecksJSON = decksJSONSerialized

	return &gameStateRow, nil
}
