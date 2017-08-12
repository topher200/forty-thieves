// Package handlers provides request handlers.
package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/topher200/forty-thieves/libdb"
)

func getCurrentUser(w http.ResponseWriter, r *http.Request) (*libdb.UserRow, bool) {
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "forty-thieves-session")
	currentUser, exists := session.Values["user"].(*libdb.UserRow)
	return currentUser, exists
}

func getIdFromPath(w http.ResponseWriter, r *http.Request) (int64, error) {
	userIdString := mux.Vars(r)["id"]
	if userIdString == "" {
		return -1, errors.New("user id cannot be empty.")
	}

	userId, err := strconv.ParseInt(userIdString, 10, 64)
	if err != nil {
		return -1, err
	}

	return userId, nil
}
