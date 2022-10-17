package handlers

import (
	"net/http"

	"example.com/program/database"
)

func RequireAuthMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !IsUserLoggedIn(r) {
			http.Redirect(w, r, "/", http.StatusUnauthorized)
			return
		}

		h(w, r)
	}
}

func RequireAdminMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUserId := CurrentUserId(r)
		user, _ := database.GetStore().UserById(*currentUserId)
		if !user.IsAdmin {
			http.Redirect(w, r, "/users", http.StatusFound)
		}

		h(w, r)
	}
}
