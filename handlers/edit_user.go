package handlers

import (
	"net/http"
	"strconv"

	"example.com/program/database"
	"example.com/program/templates"
	"example.com/program/validation"
)

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

	err = validation.ValidateUsername(username)
	if err != nil {
		templates.RenderTemplate(w, "edit", struct {
			User         *database.User
			ErrorMessage *string
		}{
			User:         user,
			ErrorMessage: ErrorMessage(err),
		})
		return
	}

	err = validation.ValidatePassword(password)
	if err != nil {
		templates.RenderTemplate(w, "edit", struct {
			User         *database.User
			ErrorMessage *string
		}{
			User:         user,
			ErrorMessage: ErrorMessage(err),
		})
		return
	}

	err = database.GetStore().EditUser(username, password, *user.Id)
	if err == nil {
		http.Redirect(w, r, "/users", http.StatusFound)
	} else {
		templates.RenderTemplate(w, "edit", struct {
			User         *database.User
			ErrorMessage *string
		}{
			User:         user,
			ErrorMessage: nil,
		})
	}

}

func ErrorMessage(err error) *string {
	if err == nil {
		return nil
	}
	errMsg := err.Error()
	return &errMsg
}
