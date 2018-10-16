/*
	File: main.go
	Contains the main program of the paraglider API.

  Information about the program:
		RESTful json API.
		Assignment 2: IGC track viewer extended - IMT2681-2018 (Cloud Technologies)

	  By Mats Ove Mandt Skjærstein.
*/

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {

	// Uses mux for regex matching on the HandleFunc paths.
	router := mux.NewRouter()

	// Functions to handle the URL paths.
	// Track:
	router.HandleFunc("/paragliding/", redirectToInfo)
	router.HandleFunc("/paragliding/api", getAPIInfo)
	router.HandleFunc("/paragliding/api/track", handleTracks)
	router.HandleFunc("/paragliding/api/track/{id:[0-9]+}", getTrackByID)
	router.HandleFunc("/paragliding/api/track/{id:[0-9]+}/{field:[a-z-A-Z-_]+}", getDetailedTrack)

	// Ticker:

	// Webhook:

	// Admin:

	// Gets the port from enviroment var.
	port := os.Getenv("PORT")
	if port == "" {
		// If there was no port, sets it to 8080.
		port = "8080"
	}

	// Starts the web-application.
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}
