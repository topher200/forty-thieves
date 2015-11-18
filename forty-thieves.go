package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/topher200/baseutil"
)

// Global var for GameState. We only support one game at a time.
var gameState GameState

func handleResources(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Providing", r.URL.Path[1:])
	http.ServeFile(w, r, r.URL.Path[1:])
}

// handleStateRequest returns a json string of the current game state.
func handleStateRequest(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(&gameState)
	baseutil.Check(err)
	fmt.Println("providing json gamestate:", string(data))
	fmt.Fprint(w, string(data))
}

func showHttp(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Root handling request:", r.URL.Path[1:])
	http.ServeFile(w, r, "res/cards.html")
}

func main() {
	gameState = NewGame()

	http.HandleFunc("/res/", handleResources)
	http.HandleFunc("/state", handleStateRequest)
	http.HandleFunc("/", showHttp)

	fmt.Println("Starting server...")
	http.ListenAndServe(":8080", nil)
}
