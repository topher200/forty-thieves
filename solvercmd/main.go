package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/topher200/forty-thieves/libdb"
	"github.com/topher200/forty-thieves/libenv"
	"github.com/topher200/forty-thieves/libgame"
	"github.com/topher200/forty-thieves/libsolver"
)

var (
	addr       = flag.String("prometheus-listen-address", ":9999", "The address to listen on for HTTP requests.")
	newGamePtr = flag.Bool(
		"new-game",
		false,
		"start a new game for analyzing. if false (default), uses latest game instead")
)

var (
	labels = []string{
		"version",
	}
	appVersion             = "v1"
	processedStatesCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "forty_thieves_processed_states_total",
			Help: "Total number of states processed",
		}, labels)
	newSavedStatesCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "forty_thieves_new_saved_states_total",
			Help: "Total number of new saved states",
		}, labels)
)

func init() {
	prometheus.MustRegister(processedStatesCounter)
	prometheus.MustRegister(newSavedStatesCounter)
}

// main process to kick off workers and solve game states
func main() {
	defer timeTrack(time.Now(), "total time")

	// host prometheus metrics
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(*addr, nil)

	// connect to database
	db, err := connectToDatabase()
	if err != nil {
		panic(fmt.Errorf("Failed to connect to database: %v.", err))
	}
	gameDB := libdb.NewGameDB(db)
	gameStateDB := libdb.NewGameStateDB(db)
	game := getOrCreateGame(gameDB, gameStateDB)

	// fire off workers
	shutdownNow := make(chan bool, 5)
	done := make(chan bool, 3)
	numWorkers := runtime.NumCPU()
	for workerId := 0; workerId < numWorkers; workerId++ {
		go doWorkerLoop(workerId, *game, shutdownNow, done)
	}
	fmt.Println("Press <enter> to exit")
	fmt.Scanln()
	fmt.Println("sending shutdown signal")
	close(shutdownNow)
	for workerId := 0; workerId < numWorkers; workerId++ {
		<-done
	}
	fmt.Println("all workers are shut down")
}

// doWorkerLoop is a helper func to pull a gameState off the queue and process it
//
// Runs until a message is seen on the 'shutdownNow' channel. Shuts itself down
// and puts a message on the 'done' channel.
func doWorkerLoop(workerId int, game libgame.Game, shutdownNow <-chan bool, done chan<- bool) {
	fmt.Printf("starting worker %d\n", workerId)

	// connect to database
	db, err := connectToDatabase()
	if err != nil {
		panic(fmt.Errorf("Failed to connect to database: %v.", err))
	}
	gameStateDB := libdb.NewGameStateDB(db)

	for {
		select {
		case _ = <-shutdownNow:
			fmt.Printf("shutting down worker %d\n", workerId)
			done <- true
			return
		default:
			// get the next game state to analyze
			gameState, err := gameStateDB.GetNextToAnalyze(game)
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
					checkGameStateSaveError(err)
				}
			} else {
				// can't flip an empty stock, nothing to do
			}

			// for each possible state we can move to, add them to the database
			for _, move := range libsolver.GetPossibleMoves(gameState) {
				if shouldSkipMove(move) {
					continue
				}

				// create a copy of our current game state
				gameStateCopy := gameState.Copy()

				// take the available move
				err = gameStateCopy.MoveCard(move)
				if err != nil {
					panic(fmt.Errorf("Error making move: %v.", err))
				}

				// save the new game state to database
				err = gameStateDB.SaveGameState(nil, gameStateCopy)
				if err == nil {
					newSavedStatesCounter.WithLabelValues(appVersion).Inc()
				} else {
					checkGameStateSaveError(err)
				}
			}

			// mark this game state as 'PROCESSED'
			err = gameStateDB.MarkAsProcessed(nil, *gameState)
			if err != nil {
				panic(fmt.Errorf("Error saving game state back to db: %v.", err))
			}
			processedStatesCounter.WithLabelValues(appVersion).Inc()
		}
	}
}

// shouldSkipMove determines whether or not we should skip a move due to business logic
//
// This function culls away moves that are likely to result in "shifting" of
// cards but not really going anywhere.
func shouldSkipMove(move libgame.MoveRequest) bool {
	if move.FromPile == libgame.FOUNDATION && move.ToPile == libgame.FOUNDATION {
		// don't keep just shifting around foundations
		return true
	}

	if move.FromPile == libgame.FOUNDATION {
		// for now, let's not let _any_ cards come down from
		// foundations. this may be changed in the future
		return true
	}

	return false // this move is fine
}

// getOrCreateGame is a helper function for getting/creating a game to process, based on user input
func getOrCreateGame(gameDB *libdb.GameDB, gameStateDB *libdb.GameStateDB) *libgame.Game {
	flag.Parse()
	var game *libgame.Game
	var err error
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
	return game
}

// connectToDatabase is a helper function to connect to our postgres db
func connectToDatabase() (db *sqlx.DB, err error) {
	dbname := "forty_thieves"
	dsn := libenv.EnvWithDefault(
		"DSN", fmt.Sprintf("postgres://postgres@localhost:5432/%s?sslmode=disable", dbname))

	db, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return db, err
}

// checkGameStateSaveError checks to the see if the given error is a dupe game. doesn't fail if it is
func checkGameStateSaveError(err error) {
	if err.Error() != "duplicate game state error" {
		panic(err)
	}
}

// timeTrack is a helper function to let us know how long something took
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
