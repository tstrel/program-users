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

	/*
		home page "/"
		- if user not logged show 2 links: loging or signup
		- if user present show message: hello Username
	*/
	// login page "/login"
	// register page "/register"

	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/register/", handlers.RegisterHandler)

	log.Println("starting server at: http://localhost")

	log.Fatal(http.ListenAndServe(":80", nil))
}
