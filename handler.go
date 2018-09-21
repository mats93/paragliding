/*
  File: handler.go
  Contains the functions that handle the API calls.
*/

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/rickb777/date/period"
)

// GET: Returns information about the API.
// Output: application/json
func getApiInfo(w http.ResponseWriter, r *http.Request) {
	// Calculates the duration since application start.
	// Uses the ISO 8601 Duration format.
	// The Date package "github.com/rickb777/date/period" is used for this.
	duration := time.Since(startTime)
	p, _ := period.NewOf(duration)
	timeStr := p.String()

	// Creates a new struct for the API info.
	currentAPI := apiInfo{timeStr, INFORMATION, VERSION}

	// Converts the struct to json.
	json, err := json.Marshal(currentAPI)
	if err != nil {
		// Sets header status code to 500 "Internal server error".
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		// Sets header content-type to application/json and status code to 200 (OK).
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Returns API info as json.
		w.Write([]byte(json))
	}
}

// POST, GET: Track registration.
// Input/Output: application/json
func handleTracks(w http.ResponseWriter, r *http.Request) {
	// Calls functions to handle the GET and POST requests.
	switch r.Method {
	case "GET":
		allTrackIDs(w, r)

	case "POST":
		insertNewTrack(w, r)

	default:
		// Sets header status code to 400 "Bad request", and writes out error.
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error: Must be a POST or GET request.")) // ToDo: Skriv ut error på en annen måte.
	}
}

// GET: Returns metadata about a given track with the provided '<id>'.
// Output: application/json
func getTrackByID(w http.ResponseWriter, r *http.Request) {
	var id int
	// Gets the ID from the URL and converts it to an integer.
	fmt.Sscanf(r.URL.Path, "/igcinfo/api/igc/%d", &id)

	// Tries to retrive the track with the requestet ID.
	// Retrieves the ID from the url.
	if rTrack, err := retriveTrackById(id); err == nil {
		// The request is valid, the track was found.

		// Converts the struct to json and outputs it.
		json, err := json.Marshal(rTrack)
		if err != nil {
			// Sets header status code to 500 "Internal server error".
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			// Sets header content-type to application/json and status code to 200 (OK).
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			// Returns the track as json.
			w.Write([]byte(json))
		}
	} else {
		// A track with the given ID does not excist.
		// Sets the header code to 404 (Not found).
		w.WriteHeader(http.StatusNotFound)
	}
}

// GET: Returns single detailed metadata about a given tracks field with the provided '<id>' and '<field>'.
// Output: text/plain.
func getDetailedTrack(w http.ResponseWriter, r *http.Request) {
	var id int
	var field string

	// Gets the field and ID from the URL, converts the ID to an integer.
	fmt.Sscanf(r.URL.Path, "/igcinfo/api/igc/%d/%s", &id, &field)

	// Tries to retrive the track with the requestet ID.
	// Retrieves the ID from the url.
	if rTrack, err := retriveTrackById(id); err == nil {
		// The request is valid, the track was found.

		// Retrieves the field specified, or 404 field not found.
		output := ""
		switch field {
		case "H_date":
			// Converts time.Time to string with string() method.
			output = rTrack.H_date.String()
		case "pilot":
			output = rTrack.Pilot
		case "glider":
			output = rTrack.Glider
		case "glider_id":
			output = rTrack.Glider_id
		case "track_length":
			// Converts float64 to string.
			output = strconv.FormatFloat(rTrack.Track_length, 'f', 6, 64)
		default:
			// If the field specified does not match any field in the track.
			// Sets the header code to 404 (Not found).
			w.WriteHeader(http.StatusNotFound)
			output = ""
		}

		// Sets header content-type to text/plain.
		w.Header().Set("Content-Type", "text/plain")

		// Returns the field.
		w.Write([]byte(output))

	} else {
		// A track with the given ID does not excist.
		// Sets the header code to 404 (Not found).
		w.WriteHeader(http.StatusNotFound)
	}
}
