package main

import (
	"log"
	"net/http"

	"example.com/program/database"
	"example.com/program/handlers"
)

func main() {
	store := database.GetStore()
	defer store.Close()

	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/register/", handlers.RegisterHandler)
	http.HandleFunc("/login/", handlers.LoginHandler)
	http.HandleFunc("/logout/", handlers.LogoutHandler)
	http.HandleFunc("/users/", handlers.RequireAuthMiddleware(handlers.ViewUsersHandler))
	http.HandleFunc("/users/delete", handlers.RequireAuthMiddleware(handlers.DeleteUserHandler))
	http.HandleFunc("/users/edit", handlers.RequireAuthMiddleware(handlers.EditUserHandler))
	log.Println("starting server at: http://localhost")

	log.Fatal(http.ListenAndServe(":80", nil))
}
