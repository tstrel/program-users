package handlers

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"example.com/program/database"
	"example.com/program/templates"
)

var currentUserId *int64

func IsUserLoggedIn() bool {
	return currentUserId != nil
}

func SetLoggedInUserID(userID int64) {
	currentUserId = &userID
}

func LogoutUser() {
	currentUserId = nil
}

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
	if IsUserLoggedIn() {
		user, err = database.GetStore().UserById(*currentUserId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	templates.RenderTemplate(w, "home", user)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if IsUserLoggedIn() {
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

	err := ValidateUsername(username)
	if err != nil {
		templates.RenderTemplate(w, "register", err)
		return
	}

	err = ValidatePassword(password)
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

	userId, err := store.CreateUser(username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	SetLoggedInUserID(userId)

	http.Redirect(w, r, "/", http.StatusFound)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	LogoutUser()
	http.Redirect(w, r, "/", http.StatusFound)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if IsUserLoggedIn() {
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

	SetLoggedInUserID(*user.Id)

	http.Redirect(w, r, "/", http.StatusFound)
}

func ViewUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, _ := database.GetStore().Users()

	// enc := json.NewEncoder(w)
	// enc.SetIndent("", "  ")
	// enc.Encode(users)   for export users

	templates.RenderTemplate(w, "users", users)
}

func ValidatePassword(password string) error {
	if len(password) <= 5 {
		return fmt.Errorf("password cannot be less than 6 characters")
	}
	if len(password) > 16 {
		return fmt.Errorf("password cannot be more than 16 characters")
	}

	if strings.Contains(password, " ") {
		return fmt.Errorf("invalid password")
	}

	return nil
}

var userNameRegExp = regexp.MustCompile("[^A-Za-z0-9]")

func ValidateUsername(username string) error {
	if len(username) <= 2 {
		return fmt.Errorf("username cannot be less than 3 characters")
	}
	if len(username) > 16 {
		return fmt.Errorf("username cannot be more than 16 characters")
	}

	if userNameRegExp.MatchString(username) {
		return fmt.Errorf("username could contain only A-Z, a-z or 0-9 characters")
	}

	return nil
}
