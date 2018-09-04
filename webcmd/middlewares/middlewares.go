// Package middlewares provides common middleware handlers.
package middlewares

import (
	"io"
	"net/http"

	"context"

	"github.com/gorilla/handlers"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
)

func SetDB(db *sqlx.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			req = req.WithContext(context.WithValue(req.Context(), "db", db))

			next.ServeHTTP(res, req)
		})
	}
}

func SetSessionStore(sessionStore sessions.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			req = req.WithContext(context.WithValue(req.Context(), "sessionStore", sessionStore))

			next.ServeHTTP(res, req)
		})
	}
}

func SetupLogger(logWriter io.Writer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return handlers.LoggingHandler(logWriter, next)
	}
}
