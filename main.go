package main

import (
	"fmt"
	"log"
	"net/http"

	"example.com/program/api"
	"example.com/program/database"
	"example.com/program/handlers"
	"github.com/gorilla/mux"
)

func main() {
	store := database.GetStore()
	defer store.Close()

	r := mux.NewRouter()

	registerApiHandlers(r)
	registerWebHandlers(r)

	log.Println("starting server at: http://localhost")
	log.Fatal(http.ListenAndServe(":80", r))
}

func registerWebHandlers(r *mux.Router) {
	r.HandleFunc("/", handlers.HomeHandler)
	r.HandleFunc("/register", handlers.RegisterHandler)
	r.HandleFunc("/login", handlers.LoginHandler)
	r.HandleFunc("/logout", handlers.LogoutHandler)

	withAuth := r.NewRoute().Subrouter()
	withAuth.Use(handlers.RequireAuthMiddleware)

	withAuth.HandleFunc("/users", handlers.ViewUsersHandler)

	withAuth.HandleFunc("/users/delete", handlers.RequireAdminMiddleware(handlers.DeleteUserHandler))

	withAuth.HandleFunc("/users/edit", handlers.RequireAdminMiddleware(handlers.EditUserHandler))
}

func registerApiHandlers(r *mux.Router) {
	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"status":"ok"}`)
	})

	apiRouter.Path("/users").Methods("GET").HandlerFunc(api.UsersHandler)              // list users
	apiRouter.Path("/users").Methods("POST").HandlerFunc(api.CreateUserHandler)        // create user
	apiRouter.Path("/users/{id}").Methods("GET").HandlerFunc(api.GetUserHandler)       // get user
	apiRouter.Path("/users/{id}").Methods("PUT").HandlerFunc(api.UpdateUserHandler)    // update user
	apiRouter.Path("/users/{id}").Methods("DELETE").HandlerFunc(api.DeleteUserHandler) // delete user
}
