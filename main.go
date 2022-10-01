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

	http.HandleFunc("/", handlers.MakeHandler(handlers.HomeHandler))
	http.HandleFunc("/register/", handlers.MakeHandler(handlers.RegisterHandler))
	// http.HandleFunc("/view/", handlers.MakeHandler(handlers.ViewHandler))
	// http.HandleFunc("/edit/", handlers.MakeHandler(handlers.EditHandler))
	// http.HandleFunc("/save/", handlers.MakeHandler(handlers.SaveHandler))

	log.Fatal(http.ListenAndServe(":80", nil))
}
