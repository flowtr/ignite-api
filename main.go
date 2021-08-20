package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func initializeRouter() {
	r := mux.NewRouter()

	r.HandleFunc("/vm", GetVMS).Methods("GET")
	r.HandleFunc("/vm/{id}", GetVM).Methods("GET")
	r.HandleFunc("/vm", CreateVM).Methods("POST")

	log.Print("Starting server on :8008")
	log.Fatal(http.ListenAndServe(":8008" +
		"", r))
}

func main() {
	initializeRouter()
}