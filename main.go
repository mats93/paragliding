/*
	File: main.go
	Contains the main program of the igc track API.

  Information about the program:
		RESTful json API.
		Assignment 1: in-memory IGC track viewer - IMT2681-2018 (Cloud Technologies)

		API to allow users to browse information about IGC files.
		The program wil not store anything in any persistant storage,
		instead it wil store submitted tracks in memory.

		This app wil run in Heroku.

		By Mats Ove Mandt Skjærstein.

  ToDo:
	  Splitt opp i forskjellige filer.
	  Skriv unit tester.
*/

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {

	// Uses mux for regex.
	router := mux.NewRouter()

	// Functions to handle the URL paths.
	router.HandleFunc("/igcinfo/api", getApiInfo)
	router.HandleFunc("/igcinfo/api/igc", handleTracks)
	router.HandleFunc("/igcinfo/api/igc/{id:[0-9]+}", getTrackByID)
	router.HandleFunc("/igcinfo/api/igc/{id:[0-9]+}/{field:[a-z-A-Z-_]+}", getDetailedTrack)

	// Gets the port from enviroment var.
	port := os.Getenv("PORT")
	if port == "" {
		// If there was no port, sets it to 8080.
		port = "8080"
	}

	// Starts the web-application.
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Panic(err)
	}
}