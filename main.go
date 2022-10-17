package main

import (
	"log"
	"net/http"

	"example.com/program/database"
	"example.com/program/handlers"
	"github.com/gorilla/mux"
)

func main() {
	store := database.GetStore()
	defer store.Close()

	r := mux.NewRouter()

	r.HandleFunc("/", handlers.HomeHandler)
	r.HandleFunc("/register", handlers.RegisterHandler)
	r.HandleFunc("/login", handlers.LoginHandler)
	r.HandleFunc("/logout", handlers.LogoutHandler)

	withAuth := r.NewRoute().Subrouter()
	withAuth.Use(handlers.RequireAuthMiddleware)

	withAuth.HandleFunc("/users", handlers.ViewUsersHandler)

	withAuth.HandleFunc("/users/delete", handlers.RequireAdminMiddleware(handlers.DeleteUserHandler))

	withAuth.HandleFunc("/users/edit", handlers.RequireAdminMiddleware(handlers.EditUserHandler))

	log.Println("starting server at: http://localhost")

	log.Fatal(http.ListenAndServe(":80", r))
}
