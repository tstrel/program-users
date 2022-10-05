package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
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

	var msg *string

	msgParam := r.URL.Query().Get("msg")
	if msgParam != "" {
		msg = &msgParam
	}

	templates.RenderTemplate(w, "users", struct {
		Users        []database.User
		ErrorMessage *string
	}{
		Users:        users,
		ErrorMessage: msg,
	})
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

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	userIdStr := q.Get("user_id")

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if userId == *currentUserId {
		msg := "You could not delete yourself"

		q := url.Values{}
		q.Add("msg", msg)

		http.Redirect(w, r, fmt.Sprintf("/users?%s", q.Encode()), http.StatusFound)
		return
	}

	if err := database.GetStore().DeleteUser(userId); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/users", http.StatusFound)
}

func EditUserHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	userIdStr := q.Get("user_id")

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	user, err := database.GetStore().UserById(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if r.Method == http.MethodGet {
		templates.RenderTemplate(w, "edit", struct {
			User         *database.User
			ErrorMessage *string
		}{
			User:         user,
			ErrorMessage: nil,
		})
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "bad request", http.StatusNotFound)
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	err = database.GetStore().EditUser(username, password, *user.Id)
	if err == nil {
		// redirect to users list
		http.Redirect(w, r, "/users", http.StatusFound)
	}

	// render form with error
}

func RequireAuthMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !IsUserLoggedIn() {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		h(w, r)
	}
}
