package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/carbocation/interpose"
	gorilla_mux "github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/topher200/forty-thieves/libenv"
	"github.com/topher200/forty-thieves/libunix"
	"github.com/topher200/forty-thieves/webcmd/handlers"
	"github.com/topher200/forty-thieves/webcmd/middlewares"
	"github.com/tylerb/graceful"
)

// NewApplication is the constructor for Application struct.
//
// If testing is true, connects to the "test" database.
func NewApplication(testing bool) (*Application, error) {
	u, err := libunix.CurrentUser()
	if err != nil {
		return nil, err
	}

	var dbname string
	if testing {
		dbname = "forty_thieves_test"
	} else {
		dbname = "forty_thieves"
	}
	dsn := libenv.EnvWithDefault(
		"DSN", fmt.Sprintf("postgres://%v@localhost:5432/%s?sslmode=disable", u, dbname))

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	cookieStoreSecret := libenv.EnvWithDefault("COOKIE_SECRET", "ittwiP92o0oi6P4i")

	app := &Application{}
	app.dsn = dsn
	app.db = db
	app.sessionStore = sessions.NewCookieStore([]byte(cookieStoreSecret))

	return app, err
}

// Application is the application object that runs HTTP server.
type Application struct {
	dsn          string
	db           *sqlx.DB
	sessionStore sessions.Store
}

func (app *Application) middlewareStruct(logWriter io.Writer) (*interpose.Middleware, error) {
	middle := interpose.New()
	middle.Use(middlewares.SetDB(app.db))
	middle.Use(middlewares.SetSessionStore(app.sessionStore))
	middle.Use(middlewares.SetupLogger(logWriter))

	middle.UseHandler(app.mux())

	return middle, nil
}

func (app *Application) mux() *gorilla_mux.Router {
	router := gorilla_mux.NewRouter()
	router.KeepContext = true

	router.HandleFunc("/", handlers.GetHome).Methods("GET").Name("/")
	router.HandleFunc("/state", handlers.HandleStateRequest)
	router.HandleFunc("/newgame", handlers.HandleNewGameRequest)
	router.HandleFunc("/move", handlers.HandleMoveRequest)
	router.HandleFunc("/flipstock", handlers.HandleFlipStockRequest)
	router.HandleFunc("/foundationcard", handlers.HandleFoundationAvailableCardRequest)

	router.PathPrefix("/bower_components").
		Handler(http.StripPrefix("/bower_components/", http.FileServer(http.Dir("bower_components")))).
		Name("/bower_components")
	router.PathPrefix("/static").Handler(http.FileServer(http.Dir(""))).Name("/static")

	return router
}

func main() {
	app, err := NewApplication(false)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	logWriter := logrus.New().Writer()
	defer logWriter.Close()
	middle, err := app.middlewareStruct(logWriter)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	serverAddress := libenv.EnvWithDefault("HTTP_ADDR", ":8888")
	certFile := libenv.EnvWithDefault("HTTP_CERT_FILE", "")
	keyFile := libenv.EnvWithDefault("HTTP_KEY_FILE", "")
	drainIntervalString := libenv.EnvWithDefault("HTTP_DRAIN_INTERVAL", "1s")

	drainInterval, err := time.ParseDuration(drainIntervalString)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	srv := &graceful.Server{
		Timeout: drainInterval,
		Server:  &http.Server{Addr: serverAddress, Handler: middle},
	}

	logrus.Infoln("Running HTTP server on " + serverAddress)
	if certFile != "" && keyFile != "" {
		err = srv.ListenAndServeTLS(certFile, keyFile)
	} else {
		err = srv.ListenAndServe()
	}
	if err != nil {
		logrus.Fatal(err.Error())
	}
}
