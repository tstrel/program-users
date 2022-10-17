package handlers

import (
	"net/http"

	"example.com/program/database"
	"example.com/program/templates"
)

func ViewUsersHandler(w http.ResponseWriter, r *http.Request) {
	currentUserId := CurrentUserId(r)
	user, _ := database.GetStore().UserById(*currentUserId)

	users, _ := database.GetStore().Users()

	var msg *string

	msgParam := r.URL.Query().Get("msg")
	if msgParam != "" {
		msg = &msgParam
	}

	templates.RenderTemplate(w, "users", struct {
		Users        []database.User
		ErrorMessage *string
		IsAdmin      bool
	}{
		Users:        users,
		ErrorMessage: msg,
		IsAdmin:      user.IsAdmin,
	})
}
