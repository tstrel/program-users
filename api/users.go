package api

import (
	"net/http"

	"example.com/program/database"
)

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	users, _ := database.GetStore().Users()

	if true {
		w.WriteHeader(http.StatusBadRequest)
		NewIndentEncoder(w).Encode(ApiError{
			Message: "bad username",
		})
	}

	apiUsers := make([]User, 0, len(users))
	for _, u := range users {
		apiUsers = append(apiUsers, User{
			Id:        u.Id,
			Username:  u.Username,
			CreatedAt: u.CreatedAt,
			IsAdmin:   u.IsAdmin,
		})
	}

	w.WriteHeader(http.StatusOK)
	NewIndentEncoder(w).Encode(UsersResponse{
		Users: apiUsers,
	})
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
}
