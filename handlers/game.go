package handlers

// handleStateRequest returns a json string of the current game state.
import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/jmoiron/sqlx"
	"github.com/topher200/deck"
	"github.com/topher200/forty-thieves/dal"
	"github.com/topher200/forty-thieves/libgame"
	"github.com/topher200/forty-thieves/libhttp"
)

// Returns the DB paramaters required to be able to get/save GameStates for this user.
func databaseParams(
	w http.ResponseWriter, r *http.Request) (*dal.GameStateDB, *dal.UserRow, error) {
	db := context.Get(r, "db").(*sqlx.DB)
	gameStateDB := dal.NewGameStateDB(db)
	currentUser, exists := getCurrentUser(w, r)
	if !exists {
		return nil, nil, errors.New("User not found")
	}
	return gameStateDB, currentUser, nil
}

func getGameState(w http.ResponseWriter, r *http.Request) (*libgame.GameState, error) {
	gameStateDB, currentUser, err := databaseParams(w, r)
	if err != nil {
		return nil, err
	}
	gameState, err := gameStateDB.GetGameState(*currentUser)
	if err != nil {
		return nil, fmt.Errorf("No game state found. %v.", err)
	}
	return gameState, nil
}

func HandleStateRequest(w http.ResponseWriter, r *http.Request) {
	gameState, err := getGameState(w, r)
	if err != nil {
		libhttp.HandleErrorJson(w, fmt.Errorf("Can't get game state: %v.", err))
		return
	}
	data, err := json.Marshal(&gameState)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	w.Header().Set("Content-Type", "text/json")
	fmt.Fprint(w, string(data))
}

// saveGameStateAndRespond saves GameState to the DB, replies with the new state.
//
// Sends a json response with the new state using the /state route.
func saveGameStateAndRespond(
	w http.ResponseWriter, r *http.Request, gameState libgame.GameState) {
	gameStateDB, currentUser, err := databaseParams(w, r)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	err = gameStateDB.SaveGameState(nil, *currentUser, gameState)
	if err != nil {
		libhttp.HandleErrorJson(w, fmt.Errorf("error saving gamestate: %v", err))
		return
	}
	HandleStateRequest(w, r)
}

// HandleNewGameRequest saves a new GameState to the DB
//
// We respond just like a /state request
func HandleNewGameRequest(w http.ResponseWriter, r *http.Request) {
	gameState := libgame.NewGame()
	saveGameStateAndRespond(w, r, gameState)
}

type MoveCommand struct {
	FromLocation string
	FromIndex    int
	ToLocation   string
	ToIndex      int
}

// HandleMoveRequest makes the move and saves the new state to the db.
//
// We respond just like a /state request
func HandleMoveRequest(w http.ResponseWriter, r *http.Request) {
	gameState, err := getGameState(w, r)
	if err != nil {
		libhttp.HandleErrorJson(w, fmt.Errorf("Can't get game state: %v.", err))
		return
	}

	// Parse the request from json
	decoder := json.NewDecoder(r.Body)
	var request MoveCommand
	err = decoder.Decode(&request)
	if err != nil {
		libhttp.HandleErrorJson(w, fmt.Errorf("failure to decode move request: %v", err))
		return
	}
	log.Printf("Handling move request from %s-%d to %s-%d\n",
		request.FromLocation, request.FromIndex, request.ToLocation, request.ToIndex)

	// Translate from string pile description to actual Decks
	parse := func(location string, index int) *deck.Deck {
		var d *deck.Deck
		switch location {
		case "tableau":
			d = &gameState.Tableaus[index]
		case "foundation":
			d = &gameState.Foundations[index]
		default:
			libhttp.HandleErrorJson(w, fmt.Errorf("unable to find deck: %v", err))
		}
		return d
	}
	from := parse(request.FromLocation, request.FromIndex)
	to := parse(request.ToLocation, request.ToIndex)

	// Move the card
	err = gameState.MoveCard(from, to)
	if err != nil {
		libhttp.HandleErrorJson(w, fmt.Errorf("invalid move: %v", err))
		return
	}

	saveGameStateAndRespond(w, r, *gameState)
}
