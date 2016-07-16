package handlers

// handleStateRequest returns a json string of the current game state.
import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/schema"
	"github.com/jmoiron/sqlx"
	"github.com/topher200/forty-thieves/dal"
	"github.com/topher200/forty-thieves/libgame"
	"github.com/topher200/forty-thieves/libhttp"
	"github.com/topher200/forty-thieves/libsolver"
)

// Returns the DB paramaters required to be able to get/save GameStates for this user.
func databaseParams(
	w http.ResponseWriter, r *http.Request) (
	*dal.GameDB, *dal.GameStateDB, *dal.UserRow, error) {
	db := context.Get(r, "db").(*sqlx.DB)
	currentUser, exists := getCurrentUser(w, r)
	if !exists {
		return nil, nil, nil, errors.New("User not found")
	}
	gameDB := dal.NewGameDB(db)
	gameStateDB := dal.NewGameStateDB(db)
	return gameDB, gameStateDB, currentUser, nil
}

func getGameState(w http.ResponseWriter, r *http.Request) (*libgame.GameState, error) {
	gameDB, gameStateDB, currentUser, err := databaseParams(w, r)
	if err != nil {
		return nil, err
	}
	gameID, err := gameDB.GetLatestGameID(*currentUser)
	if err != nil {
		return nil, err
	}
	gameState, err := gameStateDB.GetLatestGameState(gameID)
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
	_, gameStateDB, currentUser, err := databaseParams(w, r)
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

// Although it's weird, the docs want our decoder to be a global
var decoder = schema.NewDecoder()

// HandleMoveRequest makes the move and saves the new state to the db.
//
// Requests are of the form libgame.Move. The index must be included (but is
// ignored) for "stock" and "waste" piles.
//
// We respond just like a /state request
func HandleMoveRequest(w http.ResponseWriter, r *http.Request) {
	gameState, err := getGameState(w, r)
	if err != nil {
		libhttp.HandleErrorJson(w, fmt.Errorf("Can't get game state: %v.", err))
		return
	}

	// Parse the request from json
	err = r.ParseForm()
	if err != nil {
		libhttp.HandleErrorJson(
			w, fmt.Errorf("failure to decode move request: %v", err))
		return
	}
	var moveRequest libgame.MoveRequest
	err = decoder.Decode(&moveRequest, r.PostForm)
	if err != nil {
		libhttp.HandleErrorJson(
			w, fmt.Errorf("failure to decode move request: %v. form values: %v",
				err, r.PostForm))
		return
	}
	log.Printf("Handling move request from %s-%d to %s-%d\n",
		moveRequest.FromPile, moveRequest.FromIndex,
		moveRequest.ToPile, moveRequest.ToIndex)

	// Move the card
	err = gameState.MoveCard(moveRequest)
	if err != nil {
		libhttp.HandleErrorJson(w, fmt.Errorf("invalid move: %v", err))
		return
	}

	saveGameStateAndRespond(w, r, *gameState)
}

func HandleFlipStockRequest(w http.ResponseWriter, r *http.Request) {
	gameState, err := getGameState(w, r)
	if err != nil {
		libhttp.HandleErrorJson(w, fmt.Errorf("Can't get game state: %v.", err))
		return
	}

	err = gameState.FlipStock()
	if err != nil {
		libhttp.HandleErrorJson(w, fmt.Errorf("can't flip stock: %v", err))
		return
	}

	saveGameStateAndRespond(w, r, *gameState)
}

// HandleUndoMove deletes the latest move for the current user.
//
// If no error, responds with the gamestate for the new latest move (after
// deletion).
func HandleUndoMove(w http.ResponseWriter, r *http.Request) {
	// TODO(topher)
	// gameStateDB, currentUser, err := databaseParams(w, r)
	// if err != nil {
	// 	libhttp.HandleErrorJson(w, err)
	// 	return
	// }

	// err = gameStateDB.DeleteLatestGameState(nil, *currentUser)
	// if err != nil {
	// 	libhttp.HandleErrorJson(w, fmt.Errorf("Undo failed: %v", err))
	// 	return
	// }

	// HandleStateRequest(w, r)
}

func HandleFoundationAvailableCardRequest(w http.ResponseWriter, r *http.Request) {
	gameState, err := getGameState(w, r)
	if err != nil {
		libhttp.HandleErrorJson(w, fmt.Errorf("Can't get game state: %v.", err))
		return
	}

	err = libsolver.FoundationAvailableCard(gameState)
	if err != nil {
		libhttp.HandleErrorJson(w, fmt.Errorf("can't foundation card: %v", err))
		return
	}

	saveGameStateAndRespond(w, r, *gameState)
}
