package main

import (
	"container/heap"
	"fmt"

	"github.com/jinzhu/copier"
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

	for i := 0; i < 1; i++ {
		// get the next game state that is interesting
		gameState := heap.Pop(&pq).(*libgame.GameState)

		// determine the next states after ours. add them to the priority queue
		for i, move := range libsolver.GetPossibleMoves(gameState) {
			var gameStateCopy libgame.GameState
			err = copier.Copy(&gameStateCopy, &gameState)
			fmt.Println(i, move)
			if err != nil {
				panic(fmt.Errorf("Error making copy: %v.", err))
			}
			err = gameStateCopy.MoveCard(move)
			fmt.Println(gameState)
			if err != nil {
				panic(fmt.Errorf("Error making move: %v.", err))
			}
			heap.Push(&pq, &gameStateCopy)
		}

		// save the current game state to database
		err = gameStateDB.SaveGameState(nil, *gameState)
		if err != nil {
			panic(fmt.Errorf("Error saving game state to db: %v.", err))
		}
	}
}
