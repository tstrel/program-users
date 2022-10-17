package handlers

import "net/http"

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	LogoutUser()
	http.Redirect(w, r, "/", http.StatusFound)
}
