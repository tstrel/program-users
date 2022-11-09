package api

import (
	"encoding/json"
	"net/http"

	"example.com/program/database"
	"example.com/program/validation"
)

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	users, _ := database.GetStore().Users()

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
	var user UserInput
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "could not parse user input")
		return
	}
	err = validation.ValidateUsername(user.Username)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = validation.ValidatePassword(user.Password)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID, err := database.GetStore().CreateUser(user.Username, user.Password, user.IsAdmin)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "could not create a user")
		return
	}

	w.WriteHeader(http.StatusOK)
	NewIndentEncoder(w).Encode(userID)
}
