package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"example.com/program/database"
	"example.com/program/validation"
	"github.com/gorilla/mux"
)

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := (params["id"])
	userID, _ := strconv.ParseInt(ID, 10, 64)

	dbUser, _ := database.GetStore().UserById(userID)
	if dbUser == nil {
		RespondWithError(w, http.StatusNotFound, "no such user")
		return
	}

	user := User{
		Id:        dbUser.Id,
		Username:  dbUser.Username,
		CreatedAt: dbUser.CreatedAt,
		IsAdmin:   dbUser.IsAdmin,
	}

	w.WriteHeader(http.StatusOK)
	NewIndentEncoder(w).Encode(user)
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := (params["id"])
	userID, _ := strconv.ParseInt(ID, 10, 64)

	dbUser, _ := database.GetStore().UserById(userID)
	if dbUser == nil {
		RespondWithError(w, http.StatusNotFound, "no such user")
		return
	}

	var userIn UserInput
	err := json.NewDecoder(r.Body).Decode(&userIn)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "could not parse user input")
		return
	}
	err = validation.ValidateUsername(userIn.Username)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = validation.ValidatePassword(userIn.Password)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = database.GetStore().EditUser(userIn.Username, userIn.Password, userID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "server error")
		return
	}
	dbUser, _ = database.GetStore().UserById(userID)
	user := User{
		Id:        dbUser.Id,
		Username:  dbUser.Username,
		CreatedAt: dbUser.CreatedAt,
		IsAdmin:   dbUser.IsAdmin,
	}

	w.WriteHeader(http.StatusOK)
	NewIndentEncoder(w).Encode(user)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := (params["id"])
	userID, _ := strconv.ParseInt(ID, 10, 64)

	dbUser, _ := database.GetStore().UserById(userID)
	if dbUser == nil {
		RespondWithError(w, http.StatusNotFound, "no such user")
		return
	}
	err := database.GetStore().DeleteUser(userID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "server error")
		return
	}

	w.WriteHeader(http.StatusOK)
	NewIndentEncoder(w).Encode("User deleted")
}
