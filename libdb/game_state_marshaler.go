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
func UnmarshalGameState(gameStateRow GameStateRow) (gameState *libgame.GameState, err error) {
	gameState.GameID = gameStateRow.GameID
	gameState.GameStateID = gameStateRow.GameStateID
	gameState.PreviousGameState = gameStateRow.PreviousGameState
	gameState.MoveNum = gameStateRow.MoveNum

	// unmarshal json
	var deckData decksJSONStruct
	err = gameStateRow.DecksJSON.Unmarshal(deckData)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling gameStateRow: %v", err)
	}
	gameState.Stock = deckData.Stock
	gameState.Foundations = deckData.Foundations
	gameState.Tableaus = deckData.Tableaus
	gameState.Waste = deckData.Waste

	return gameState, nil
}

func MarshalGameState(gameState libgame.GameState) (gameStateRow *GameStateRow, err error) {
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

	panic(decksJSON)

	dataStruct := GameStateRow{}
	dataStruct.GameID = gameState.GameID
	dataStruct.GameStateID = gameState.GameStateID
	dataStruct.MoveNum = gameState.MoveNum
	dataStruct.PreviousGameState = gameState.PreviousGameState
	dataStruct.Score = gameState.Score
	dataStruct.DecksJSON = decksJSONSerialized

	return gameStateRow, nil
}
