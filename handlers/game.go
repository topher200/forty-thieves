package handlers

// handleStateRequest returns a json string of the current game state.
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/topher200/baseutil"
	"github.com/topher200/deck"
	"github.com/topher200/forty-thieves/libgame"
)

func HandleStateRequest(w http.ResponseWriter, r *http.Request) {
	// temp
	gameState := libgame.NewGame()

	w.Header().Set("Content-Type", "text/json")
	data, err := json.Marshal(&gameState)
	baseutil.Check(err)
	log.Println("providing json gamestate:", string(data))
	fmt.Fprint(w, string(data))
}

func HandleMoveRequest(w http.ResponseWriter, r *http.Request) {
	// temp
	gameState := libgame.NewGame()

	// Parse the request from json
	type Message struct {
		FromLocation string
		FromIndex    int
		ToLocation   string
		ToIndex      int
	}
	decoder := json.NewDecoder(r.Body)
	var request Message
	err := decoder.Decode(&request)
	if err != nil {
		log.Println("Failed to decode move request:", r.Body)
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
			panic("Unable to find deck")
		}
		return d
	}
	from := parse(request.FromLocation, request.FromIndex)
	to := parse(request.ToLocation, request.ToIndex)

	// Move the card
	err = gameState.MoveCard(from, to)
	if err == nil {
		http.Redirect(w, r, "/", 200)
	} else {
		http.Error(w, err.Error(), 400)
	}
}