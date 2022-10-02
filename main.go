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

	log.Println("starting server at: http://localhost")

	log.Fatal(http.ListenAndServe(":80", nil))
}
