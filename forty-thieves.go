package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/topher200/baseutil"
	"github.com/topher200/deck"
)

// Global var for GameState. We only support one game at a time.
var gameState GameState

func handleResources(w http.ResponseWriter, r *http.Request) {
	log.Println("Providing res", r.URL.Path)
	http.ServeFile(w, r, r.URL.Path[1:])
}

func handleJavascript(w http.ResponseWriter, r *http.Request) {
	log.Println("Providing javascript", r.URL.Path)
	w.Header().Set("Content-Type", "application/javascript")
	http.ServeFile(w, r, r.URL.Path[1:])
}

// handleStateRequest returns a json string of the current game state.
func handleStateRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/json")
	data, err := json.Marshal(&gameState)
	baseutil.Check(err)
	log.Println("providing json gamestate:", string(data))
	fmt.Fprint(w, string(data))
}

func handleMoveRequest(w http.ResponseWriter, r *http.Request) {
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

func showHttp(w http.ResponseWriter, r *http.Request) {
	log.Println("Root handling request:", r.URL)
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, "res/game.html")
}

func parseCommandLine() {
	deterministicPointer := kingpin.Flag("deterministic",
		"makes our output deterministic by allowing the default rand.Seed").
		Short('d').Bool()
	kingpin.Parse()

	if !*deterministicPointer {
		log.Println("Seeded randomly")
		rand.Seed(time.Now().UTC().UnixNano())
	} else {
		log.Println("Seeded deterministically")
	}
}

func main() {
	parseCommandLine()

	gameState = NewGame()

	http.HandleFunc("/res/", handleResources)
	http.HandleFunc("/bower_components/", handleJavascript)
	http.HandleFunc("/state", handleStateRequest)
	http.HandleFunc("/move", handleMoveRequest)
	http.HandleFunc("/", showHttp)

	log.Println("Starting server...")
	http.ListenAndServe(":8080", nil)
}
