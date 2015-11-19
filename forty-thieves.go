package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/topher200/baseutil"
	"github.com/topher200/deck"
)

// Global var for GameState. We only support one game at a time.
var gameState GameState

func handleResources(w http.ResponseWriter, r *http.Request) {
	log.Println("Providing", r.URL.Path[1:])
	http.ServeFile(w, r, r.URL.Path[1:])
}

// handleStateRequest returns a json string of the current game state.
func handleStateRequest(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(&gameState)
	baseutil.Check(err)
	log.Println("providing json gamestate:", string(data))
	fmt.Fprint(w, string(data))
}

func handleMoveRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling move request", r)

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
		log.Println("failed to decode move request:", r.Body)
		return
	}
	log.Printf("Moving from %s-%d to %s-%d\n",
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
	gameState.MoveCard(from, to)
}

func showHttp(w http.ResponseWriter, r *http.Request) {
	log.Println("Root handling request:", r.URL.Path[1:])
	http.ServeFile(w, r, "res/cards.html")
}

func main() {
	gameState = NewGame()

	http.HandleFunc("/res/", handleResources)
	http.HandleFunc("/state", handleStateRequest)
	http.HandleFunc("/move", handleMoveRequest)
	http.HandleFunc("/", showHttp)

	log.Println("Starting server...")
	http.ListenAndServe(":8080", nil)
}
