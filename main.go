package main

import (
	"encoding/gob"
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
	"github.com/topher200/forty-thieves/dal"
	"github.com/topher200/forty-thieves/handlers"
	"github.com/topher200/forty-thieves/libenv"
	"github.com/topher200/forty-thieves/libunix"
	"github.com/topher200/forty-thieves/middlewares"
	"github.com/tylerb/graceful"
)

func init() {
	gob.Register(&dal.UserRow{})
}

// NewApplication is the constructor for Application struct.
//
// If testing is true, connects to the "test" database.
func NewApplication(testing bool) (*Application, error) {
	u, err := libunix.CurrentUser()
	if err != nil {
		return nil, err
	}

	dbname := "forty-thieves"
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
	MustLogin := middlewares.MustLogin

	router := gorilla_mux.NewRouter()
	router.KeepContext = true

	router.Handle("/", MustLogin(http.HandlerFunc(handlers.GetHome))).Methods("GET").Name("/")

	router.HandleFunc("/signup", handlers.GetSignup).Methods("GET").Name("/signup.Get")
	router.HandleFunc("/signup", handlers.PostSignup).Methods("POST").Name("/signup.Post")
	router.HandleFunc("/login", handlers.GetLogin).Methods("GET").Name("/login.Get")
	router.HandleFunc("/login", handlers.PostLogin).Methods("POST").Name("/login.Post")
	router.HandleFunc("/logout", handlers.GetLogout).Methods("GET").Name("/logout.Get")

	router.Handle(
		"/users/{id:[0-9]+}",
		MustLogin(http.HandlerFunc(handlers.PostPutDeleteUsersID))).
		Methods("POST", "PUT", "DELETE").
		Name("/users/{id}")

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
