package handlers

import (
	"net/http"

	"example.com/program/database"
	"example.com/program/templates"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if IsUserLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	if r.Method == http.MethodGet {
		templates.RenderTemplate(w, "login", nil)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "bad request", http.StatusNotFound)
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	store := database.GetStore()

	user, _ := store.UserByName(username)
	if user == nil {
		templates.RenderTemplate(w, "login", "no such user")
		return
	}

	if password != user.Password {
		templates.RenderTemplate(w, "login", "password does not match")
		return
	}

	SetLoggedInUserID(w, r, *user.Id)

	http.Redirect(w, r, "/", http.StatusFound)
}
