package api

import (
	"time"
)

type User struct {
	Id        *int64    `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	IsAdmin   bool      `json:"is_admin"`
}

type UserInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

type UsersResponse struct {
	Users []User `json:"users"`
}

type ApiError struct {
	Message string `json:"message"`
}
