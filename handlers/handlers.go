package handlers

import (
	"net/http"

	"example.com/program/database"
	"example.com/program/templates"
)

var currentUserId *int64

func HomeHandler(w http.ResponseWriter, r *http.Request) {
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
	templates.RenderTemplate(w, "register", nil)
}

func SaveRegistration(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	userId, err := database.GetStore().CreateUser(username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	currentUserId = &userId

	http.Redirect(w, r, "/", http.StatusFound)
}
