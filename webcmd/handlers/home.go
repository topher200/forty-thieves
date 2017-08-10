package handlers

import (
	"html/template"
	"net/http"

	"github.com/topher200/forty-thieves/dal"
	"github.com/topher200/forty-thieves/libhttp"
)

func GetHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	currentUser, exists := getCurrentUser(w, r)
	if !exists {
		http.Redirect(w, r, "/logout", 302)
		return
	}

	data := struct {
		CurrentUser *dal.UserRow
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
