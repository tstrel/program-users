package handlers

import (
	"net/http"

	"example.com/program/database"
	"example.com/program/templates"
	"example.com/program/validation"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if IsUserLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	if r.Method == http.MethodGet {
		templates.RenderTemplate(w, "register", nil)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "bad request", http.StatusNotFound)
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	adminRight := r.FormValue("adminRight")

	var isAdmin bool
	if adminRight == "is_admin" {
		isAdmin = true
	}

	err := validation.ValidateUsername(username)
	if err != nil {
		templates.RenderTemplate(w, "register", err)
		return
	}

	err = validation.ValidatePassword(password)
	if err != nil {
		templates.RenderTemplate(w, "register", err)
		return
	}

	store := database.GetStore()

	user, _ := store.UserByName(username)
	if user != nil {
		templates.RenderTemplate(w, "register", "such user already exists")
		return
	}

	userId, err := store.CreateUser(username, password, isAdmin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	SetLoggedInUserID(w, r, userId)

	http.Redirect(w, r, "/", http.StatusFound)
}
