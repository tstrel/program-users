package handlers

import "net/http"

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	LogoutUser(w, r)
	http.Redirect(w, r, "/", http.StatusFound)
}
