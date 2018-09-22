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
	firstGameState := libgame.DealNewGame(*game)
	err = gameStateDB.SaveGameState(nil, firstGameState)
	if err != nil {
		panic(fmt.Errorf("Error saving new game's first gamestate: %v.", err))
	}

	pq := make(PriorityQueue, 1)
	pq[0] = &firstGameState

	for pq.Len() > 0 {
		// get the next game state that is interesting
		gameState := heap.Pop(&pq).(*libgame.GameState)

		// determine the next states after ours. add them to the priority queue
		for _, move := range libsolver.GetPossibleMoves(gameState) {
			// create a copy of our current game state
			gameStateCopy := gameState.Copy()
			if err != nil {
				panic(fmt.Errorf("Error making copy: %v.", err))
			}

			// take the available move
			err = gameStateCopy.MoveCard(move)
			if err != nil {
				panic(fmt.Errorf("Error making move: %v.", err))
			}

			// save the new game state to database
			err = gameStateDB.SaveGameState(nil, gameStateCopy)
			if err != nil {
				panic(fmt.Errorf("Error saving game state to db: %v.", err))
			}

			// push the new game state onto the heap to be processed
			heap.Push(&pq, &gameStateCopy)
		}
	}
}
