// Package middlewares provides common middleware handlers.
package middlewares

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
)

func SetDB(db *sqlx.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			context.Set(req, "db", db)

			next.ServeHTTP(res, req)
		})
	}
}

func SetCookieStore(cookieStore *sessions.CookieStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			context.Set(req, "cookieStore", cookieStore)

			next.ServeHTTP(res, req)
		})
	}
}

// MustLogin is a middleware that checks existence of current user.
func MustLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cookieStore := context.Get(req, "cookieStore").(*sessions.CookieStore)
		session, _ := cookieStore.Get(req, "forty-thieves-session")
		userRowInterface := session.Values["user"]

		if userRowInterface == nil {
			http.Redirect(res, req, "/login", 302)
			return
		}

		next.ServeHTTP(res, req)
	})
}

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		next.ServeHTTP(res, req)
		routeUsed := mux.CurrentRoute(req)
		var routeName string
		if routeUsed != nil {
			routeName = routeUsed.GetName()
		}
		logrus.Infof("req: %v\n  handled by: %v\n  resp: %v", req, routeName, res)
	})
}
