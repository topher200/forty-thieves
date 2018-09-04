package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/topher200/forty-thieves/libdb"
	"github.com/topher200/forty-thieves/libenv"
	"github.com/topher200/forty-thieves/libgame"
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
		//   request a move state that is interesting
		//   determine the next state for that game
		//   save to database
	}
}
