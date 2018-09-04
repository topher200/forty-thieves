package main

import (
	"container/heap"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/topher200/forty-thieves/libdb"
	"github.com/topher200/forty-thieves/libenv"
	"github.com/topher200/forty-thieves/libgame"
	"github.com/topher200/forty-thieves/libsolver"
)

func ConnectToDatabase() (db *sqlx.DB, err error) {
	dbname := "forty_thieves"
	dsn := libenv.EnvWithDefault(
		"DSN", fmt.Sprintf("postgres://postgres@localhost:5432/%s?sslmode=disable", dbname))

	db, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return db, err
}

func main() {
	// connect to database
	db, err := ConnectToDatabase()
	gameDB := libdb.NewGameDB(db)
	gameStateDB := libdb.NewGameStateDB(db)

	// create a game
	game, err := gameDB.CreateNewGame(nil)
	if err != nil {
		panic(fmt.Errorf("Error creating new game: %v.", err))
	}
	libgame.DealNewGame(*game)

	pq := make(PriorityQueue, 1)
	pq[0], err = gameStateDB.GetFirstGameState(*game)
	if err != nil {
		panic(fmt.Errorf("Error getting new game gamestate: %v.", err))
	}

	for i := 0; i < 100; i++ {
		// get a move state that is interesting
		gameState := heap.Pop(&pq)

		// determine the next states for that state. add them to the priority queue
		for _, move := range libsolver.GetPossibleMoves(*gameState) {
			heap.Push(&pq, move)
		}

		// save this game state to database
		err = gameStateDB.SaveGameState(nil, gameState)
		if err != nil {
			panic(fmt.Errorf("Error saving game state to db: %v.", err))
		}
	}
}
