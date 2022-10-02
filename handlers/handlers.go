package handlers

import (
	"fmt"
	"net/http"

	"example.com/program/database"
	"example.com/program/templates"
)

var currentUserId *int64

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Not found")
		return
	}

	var (
		user *database.User
		err  error
	)
	if currentUserId != nil {
		user, err = database.GetStore().UserById(*currentUserId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	templates.RenderTemplate(w, "home", user)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		templates.RenderTemplate(w, "register", nil)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "bad request", http.StatusNotFound)
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	store := database.GetStore()

	// u, err := store.UserByName(1)
	if true { // if user exists
		templates.RenderTemplate(w, "register", "such user already exists")
		return
	}

	userId, err := store.CreateUser(username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	currentUserId = &userId

	http.Redirect(w, r, "/", http.StatusFound)
}
