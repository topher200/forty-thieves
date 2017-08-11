package handlers

// handleStateRequest returns a json string of the current game state.
import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/schema"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"github.com/topher200/deck"
	"github.com/topher200/forty-thieves/libgame"
	"github.com/topher200/forty-thieves/libhttp"
	"github.com/topher200/forty-thieves/libsolver"
	"github.com/topher200/forty-thieves/libdb"
)

// Returns the DB paramaters required to be able to get/save GameStates for this user.
func databaseParams(
	w http.ResponseWriter, r *http.Request) (*dal.GameDB, *dal.GameStateDB, *dal.UserRow, error) {
	db := r.Context().Value("db").(*sqlx.DB)
	gameDB := dal.NewGameDB(db)
	gameStateDB := dal.NewGameStateDB(db)
	currentUserRow, exists := getCurrentUser(w, r)
	if !exists {
		return nil, nil, nil, errors.New("User not found")
	}
	return gameDB, gameStateDB, currentUserRow, nil
}

// parseGameStateIdFromQuery gets the game state id from the URL
//
// Returns error on unexpected error. Returns uuid.Nil if no UUID is found.
func parseGameStateIdFromQuery(r *http.Request) (uuid.UUID, error) {
	queryStringValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return uuid.Nil, err
	}

	gameStateIDString := queryStringValues.Get("gameStateID")
	if queryStringValues.Get("gameStateID") != "" {
		gameStateID, err := uuid.FromString(gameStateIDString)
		if err != nil {
			return uuid.Nil, err
		}
		return gameStateID, nil
	} else {
		// request is empty
		return uuid.Nil, nil
	}
}

// parseGameStateFromQuery parses the gameStateID and returns the GameState
func parseGameStateFromQuery(w http.ResponseWriter, r *http.Request) (*libgame.GameState, error) {
	// check query param
	gameStateID, err := parseGameStateIdFromQuery(r)
	if err != nil {
		return nil, err
	}
	if gameStateID == uuid.Nil {
		return nil, fmt.Errorf("Game state id required in query param")
	}

	// get the referenced game state
	_, gameStateDB, _, err := databaseParams(w, r)
	if err != nil {
		return nil, err
	}
	// NOTE: we currently don't do any checking to make sure game state id
	// and user id match
	gameState, err := gameStateDB.GetGameStateById(gameStateID)
	if err != nil {
		return nil, fmt.Errorf("Game state id %v not found: %v", gameStateID, err)
	}
	return gameState, nil
}

// replyWithGameState sends a JSON reponse with the given game state
func replyWithGameState(w http.ResponseWriter, r *http.Request, gameState libgame.GameState) {
	type GameStateWithChildren struct {
		GameID            int64
		GameStateID       uuid.UUID
		PreviousGameState uuid.NullUUID
		MoveNum           int64
		Stock             deck.Deck
		Foundations       []deck.Deck
		Tableaus          []deck.Deck
		Waste             deck.Deck
		Score             int
		ChildGameStates   []uuid.UUID
	}

	// get children
	_, gameStateDB, _, err := databaseParams(w, r)
	childGameStates, err := gameStateDB.GetChildGameStates(gameState)
	if err != nil {
		libhttp.HandleServerError(w, err)
		return
	}

	// make new struct with children
	gs := GameStateWithChildren{
		gameState.GameID,
		gameState.GameStateID,
		gameState.PreviousGameState,
		gameState.MoveNum,
		gameState.Stock,
		gameState.Foundations,
		gameState.Tableaus,
		gameState.Waste,
		gameState.Score,
		childGameStates,
	}

	// convert to json and send
	data, err := json.Marshal(&gs)
	if err != nil {
		libhttp.HandleServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "text/json")
	fmt.Fprint(w, string(data))
}

// TODO(topher): describe response
func HandleStateRequest(w http.ResponseWriter, r *http.Request) {
	gameStateID, err := parseGameStateIdFromQuery(r)
	if err != nil {
		libhttp.HandleServerError(w, err)
		return
	}

	gameDB, gameStateDB, currentUserRow, err := databaseParams(w, r)
	if err != nil {
		libhttp.HandleServerError(w, err)
		return
	}
	// If the game state id for the request is empty, send the latest game
	// state for that user. Otherwise, send the one requested
	var gameState *libgame.GameState
	if gameStateID != uuid.Nil {
		// NOTE: if the user provides a game state id, we currently
		// don't do any checking against the user id to make sure they
		// match
		logrus.Infof("getting gamestate for gamestate id %v", gameStateID)
		gameState, err = gameStateDB.GetGameStateById(gameStateID)
		if err != nil {
			libhttp.HandleClientError(w,
				fmt.Errorf("Game state id %v not found: %v", gameStateID, err),
				http.StatusBadRequest)
			return
		}
	} else {
		logrus.Infof("no game state provided - finding latest game")
		// request is empty, find the latest
		game, err := gameDB.GetLatestGame(*currentUserRow)
		if err != nil {
			libhttp.HandleClientError(w, fmt.Errorf("No game found: %v", err), http.StatusBadRequest)
			return
		}
		gameState, err = gameStateDB.GetFirstGameState(*game)
		if err != nil {
			libhttp.HandleServerError(w, err)
			return
		}
	}

	replyWithGameState(w, r, *gameState)
}

// saveGameStateAndRespond saves GameState to the DB, replies with the new state.
//
// Sends a json response with the new state using the /state route.
func saveGameStateAndRespond(
	w http.ResponseWriter, r *http.Request, gameState libgame.GameState) {
	_, gameStateDB, _, err := databaseParams(w, r)
	if err != nil {
		libhttp.HandleServerError(w, fmt.Errorf("Error getting database params: %v.", err))
		return
	}
	err = gameStateDB.SaveGameState(nil, gameState)
	if err != nil {
		libhttp.HandleServerError(w, fmt.Errorf("error saving gamestate: %v", err))
		return
	}
	replyWithGameState(w, r, gameState)
}

// HandleNewGameRequest saves a new GameState to the DB
//
// We respond just like a /state request
func HandleNewGameRequest(w http.ResponseWriter, r *http.Request) {
	gameDB, _, currentUserRow, err := databaseParams(w, r)
	if err != nil {
		libhttp.HandleServerError(w, fmt.Errorf("Error getting database params: %v.", err))
		return
	}
	game, err := gameDB.CreateNewGame(nil, *currentUserRow)
	if err != nil {
		libhttp.HandleServerError(w, fmt.Errorf("Error creating new game: %v.", err))
		return
	}
	gameState := libgame.DealNewGame(*game)
	saveGameStateAndRespond(w, r, gameState)
}

// Although it's weird, the docs want our decoder to be a global
var decoder = schema.NewDecoder()

// HandleMoveRequest makes the move and saves the new state to the db.
//
// Requests are of the form libgame.Move. The index must be included (but is
// ignored) for "stock" and "waste" piles.
//
// TODO(topher): change these to "respond like HandleMoveRequest"
// We respond just like a /state request
func HandleMoveRequest(w http.ResponseWriter, r *http.Request) {
	gameState, err := parseGameStateFromQuery(w, r)
	if err != nil {
		libhttp.HandleServerError(w, fmt.Errorf("failure to get game state: %v", err))
		return
	}

	// Parse the request from json
	err = r.ParseForm()
	if err != nil {
		libhttp.HandleServerError(
			w, fmt.Errorf("failure to decode move request: %v", err))
		return
	}
	var moveRequest libgame.MoveRequest
	err = decoder.Decode(&moveRequest, r.PostForm)
	if err != nil {
		libhttp.HandleServerError(
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
		libhttp.HandleClientError(w, fmt.Errorf("invalid move: %v", err), http.StatusBadRequest)
		return
	}

	saveGameStateAndRespond(w, r, *gameState)
}

func HandleFlipStockRequest(w http.ResponseWriter, r *http.Request) {
	gameState, err := parseGameStateFromQuery(w, r)
	if err != nil {
		libhttp.HandleServerError(w, fmt.Errorf("failure to get game state: %v", err))
		return
	}

	err = gameState.FlipStock()
	if err != nil {
		libhttp.HandleClientError(w, fmt.Errorf("can't flip stock: %v", err),
			http.StatusBadRequest)
		return
	}

	saveGameStateAndRespond(w, r, *gameState)
}

func HandleFoundationAvailableCardRequest(w http.ResponseWriter, r *http.Request) {
	gameState, err := parseGameStateFromQuery(w, r)
	if err != nil {
		libhttp.HandleServerError(w, fmt.Errorf("failure to get game state: %v", err))
		return
	}

	err = libsolver.FoundationAvailableCard(gameState)
	if err != nil {
		libhttp.HandleClientError(w, fmt.Errorf("can't foundation any cards: %v", err),
			http.StatusBadRequest)
		return
	}

	saveGameStateAndRespond(w, r, *gameState)
}
