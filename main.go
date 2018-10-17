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
	router.HandleFunc("/paragliding/", RedirectToInfo)
	router.HandleFunc("/paragliding/api", GetAPIInfo)
	router.HandleFunc("/paragliding/api/track", HandleTracks)
	router.HandleFunc("/paragliding/api/track/{id:[0-9]+}", GetTrackByID)
	router.HandleFunc("/paragliding/api/track/{id:[0-9]+}/{field:[a-z-A-Z-_]+}", GetDetailedTrack)

	/* Ticker:
	router.HandleFunc("/paragliding/api/ticker/latest", )
	router.HandleFunc("/paragliding/api/ticker/", )
	router.HandleFunc("/paragliding/api/ticker/<timestamp>", )

	// Webhook:
	router.HandleFunc("/paragliding/api/webhook/new_track/", ) // POST
	router.HandleFunc("/paragliding/api/webhook/new_track/<webhook_id>", ) GET og DELETE

	// Admin:
	router.HandleFunc("/paragliding/admin/api/tracks_count", )
	router.HandleFunc("/paragliding/admin/api/tracks", ) // DELETE
	*/

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
