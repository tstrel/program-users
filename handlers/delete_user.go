package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"example.com/program/database"
)

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
