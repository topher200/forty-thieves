package handlers

import (
	"html/template"
	"net/http"

	"github.com/topher200/forty-thieves/libhttp"
	"github.com/topher200/forty-thieves/libdb"
)

func GetHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	currentUser, exists := getCurrentUser(w, r)
	if !exists {
		http.Redirect(w, r, "/logout", 302)
		return
	}

	data := struct {
		CurrentUser *libdb.UserRow
	}{
		currentUser,
	}

	tmpl, err := template.ParseFiles("templates/dashboard.html.tmpl", "templates/home.html.tmpl")
	if err != nil {
		libhttp.HandleServerError(w, err)
		return
	}

	tmpl.Execute(w, data)
}
