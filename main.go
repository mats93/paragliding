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
	"time"

	"github.com/gorilla/mux"
	"github.com/mats93/paragliding/admin"
	"github.com/mats93/paragliding/ticker"
	"github.com/mats93/paragliding/track"
)

// COLLECTION is the collection to be used when the main app is running.
const COLLECTION = "Tracks"

// The start time for the API service.
var startTime = time.Now()

func main() {

	// Injects the startime to the track package.
	track.StartTime = startTime

	// Injects the MongoDB collection to use.
	track.Collection = COLLECTION
	admin.Collection = COLLECTION
	ticker.Collection = COLLECTION

	// Uses mux for regex matching on the HandleFunc paths.
	router := mux.NewRouter()

	// Functions to handle the URL paths.
	// Track:
	router.HandleFunc("/paragliding/", track.RedirectToInfo)
	router.HandleFunc("/paragliding/api", track.GetAPIInfo)
	router.HandleFunc("/paragliding/api/track", track.HandleTracks)
	router.HandleFunc("/paragliding/api/track/{id:[0-9]+}", track.GetTrackByID)
	router.HandleFunc("/paragliding/api/track/{id:[0-9]+}/{field:[a-z-A-Z-_]+}", track.GetDetailedTrack)

	// Ticker:
	router.HandleFunc("/paragliding/api/ticker/latest", ticker.GetLastTimestamp)
	router.HandleFunc("/paragliding/api/ticker/", ticker.GetTimestamps)
	router.HandleFunc("/paragliding/api/ticker/{timestamp:[0-9]+}", ticker.GetTimestampsNewerThen)

	/* Webhook:
	router.HandleFunc("/paragliding/api/webhook/new_track/", ) // POST
	router.HandleFunc("/paragliding/api/webhook/new_track/<webhook_id>", ) GET og DELETE
	*/

	// Admin:
	router.HandleFunc("/paragliding/admin/api/tracks_count", admin.GetTrackCount)
	router.HandleFunc("/paragliding/admin/api/tracks", admin.DeleteAllTracks)
	router.HandleFunc("/paragliding/admin/api/insert", admin.InsertSomeTracks) // ToDo: Remove this.

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
