package main

import (
	"flag"
	"fmt"
	"log"
	"time"

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
	defer timeTrack(time.Now(), "total time")
	// connect to database
	db, err := ConnectToDatabase()
	gameDB := libdb.NewGameDB(db)
	gameStateDB := libdb.NewGameStateDB(db)

	newGamePtr := flag.Bool(
		"new-game",
		false,
		"start a new game for analyzing. if false (default), uses latest game instead")
	flag.Parse()
	var game *libgame.Game
	if *newGamePtr {
		// create a game
		game, err = gameDB.CreateNewGame(nil)
		if err != nil {
			panic(fmt.Errorf("Error creating new game: %v.", err))
		}
		firstGameState := libgame.DealNewGame(*game)
		err = gameStateDB.SaveGameState(nil, firstGameState)
		if err != nil {
			panic(fmt.Errorf("Error saving new game's first gamestate: %v.", err))
		}
	} else {
		// use existing game
		game, err = gameDB.GetLatestGame()
		if err != nil {
			panic(fmt.Errorf("Error getting game: %v.", err))
		}
	}

	defer timeTrack(time.Now(), "processing loop")
	for true {
		// get the next game state to analyze
		gameState, err := gameStateDB.GetNextToAnalyze(*game)
		if err != nil {
			panic(fmt.Errorf("Error getting next game state to analyze: %v.", err))
		}

		// if the best game state is solved, we're done!
		if gameState.Score == 0 {
			break
		}

		// flip the stock and save that new state to the database
		gameStateCopy := gameState.Copy()
		err = gameStateCopy.FlipStock()
		if err == nil {
			err = gameStateDB.SaveGameState(nil, gameStateCopy)
			if err != nil {
				panic(fmt.Errorf("Error saving flipped game state to db: %v.", err))
			}
		} else {
			// can't flip an empty stock, nothing to do
		}

		// for each possible state we can move to, add them to the database
		for _, move := range libsolver.GetPossibleMoves(gameState) {
			// create a copy of our current game state
			gameStateCopy := gameState.Copy()

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
		}

		// mark this game state as 'PROCESSED'
		err = gameStateDB.MarkAsProcessed(nil, *gameState)
		if err != nil {
			panic(fmt.Errorf("Error saving game state back to db: %v.", err))
		}
	}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
