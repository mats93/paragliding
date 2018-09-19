/*
RESTful json API.
Assignment 1: in-memory IGC track viewer - IMT2681-2018 (Cloud Technologies)

API to allow users to browse information about IGC files.
The program wil not store anything in any persistant storage,
instead it wil store submitted tracks in memory.

This app wil run in Heroku.

ToDo:
  Add "content-type" for header.
  Sjekk om response header (404 osv) stemmer.
	Splitt opp i forskjellige filer.
	Skriv unit tester.


By Mats
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	igc "github.com/marni/goigc"
	"github.com/rickb777/date/period"
)

// ***************************************************
// Constants:

// Const for the API information.
const INFORMATION = "Service for IGC tracks"
const VERSION = "v1"

// ***************************************************
// Structs:

// Format for the IGC file.
type track struct {
	Id           int       `json:"-"`
	H_date       time.Time `json:"H_date"`
	Pilot        string    `json:"pilot"`
	Glider       string    `json:"glider"`
	Glider_id    string    `json:"glider_id"`
	Track_length float64   `json:"track_length"`
}

// Format for the API information.
type apiInfo struct {
	Uptime  string `json:"uptime"`
	Info    string `json:"info"`
	Version string `json:"version"`
}

// Format for the url information.
type url struct {
	Url string `json:"url"`
}

// ***************************************************
// Global variables:

// Slice with structs of type "track" (in-memory DB).
var trackSlice []track

// Generates an unique ID for all tracks.
var lastUsedID int = 0

// The start time for the API service.
var startTime time.Time = time.Now()

// ***************************************************
// Functions:

// Searches for a given track with an ID.
func retriveTrackById(id int) (track, error) {
	// Loops through all tracks.
	for _, tr := range trackSlice {
		// If a track with the same ID as the parameter is found, return the track and no error.
		if tr.Id == id {
			return tr, nil
		}
	}
	// If no tracks was found, return an empty track and an error.
	emptyTrack := track{}
	return emptyTrack, errors.New("Could not find any track with given ID")
}

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

// GET: Returns an array of al track IDs.
// Output: application/json.
func allTrackIDs(w http.ResponseWriter, r *http.Request) {
	// Converts the IDs to a "json list".
	// Loops through all IDs and formats them.
	message := "["
	for i := 0; i < len(trackSlice); i++ {
		// Converts int to string.
		message += strconv.Itoa(trackSlice[i].Id)
		if i != len(trackSlice)-1 {
			message += ", "
		}
	}
	// The end for the "json list".
	message += "]"

	// Sets header content-type to application/json and status code to 200 (OK).
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Outputs the list (empty or full).
	w.Write([]byte(string(message)))
}

// POST: Takes the post request as json format and inserts a new track to the "DB".
// Input/Output: application/json
func insertNewTrack(w http.ResponseWriter, r *http.Request) {
	var newUrl url

	// Decodes the json url and converts it to a struct.
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newUrl)

	// If the decoding fails:
	if err != nil {
		// Sets header status code to 400 "Bad request".
		w.WriteHeader(http.StatusBadRequest)

		// If the json decoding works:
	} else {
		// Adds the new track from the URL.
		trackFile, err := igc.ParseLocation(newUrl.Url)

		// If the igc parser fails:
		if err != nil {
			// Sets header status code to 400 "Bad request".
			w.WriteHeader(http.StatusBadRequest)

			// If the parser gets the IGC content:
		} else {
			// Calculates the total distance for the track.
			var sum float64
			// Loops through all Points[] in the track.
			for i := 0; i < len(trackFile.Points)-1; i++ {
				// Adds the distance between two points togheter.
				sum += trackFile.Points[i].Distance(trackFile.Points[i+1])
			}

			// Increments the global ID, to be used for an internal track ID.
			lastUsedID++

			// Adds the new track to a slice (memory-DB).
			trackSlice = append(trackSlice,
				track{lastUsedID, trackFile.Header.Date, trackFile.Pilot, trackFile.GliderType,
					trackFile.GliderID, sum})

			// Converts the id to json format.
			jsonID := fmt.Sprintf("{\"id\": %d}", lastUsedID)

			// Sets header content-type to application/json and status code to 200 (OK).
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			// Returns the given tracks ID as json.
			w.Write([]byte(jsonID))
		}
	}
}

// POST, GET: Track registration.
// Input/Output: application/json
func handleTracks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	// Calls functions to handle the GET and POST requests.
	case "GET":
		allTrackIDs(w, r)

	case "POST":
		insertNewTrack(w, r)

	default:
		// Sets header status code to 400 "Bad request", and writes out error.
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error: Must be a POST or GET request."))
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

	/*
		Id           int       `json:"-"`
		H_date       time.Time `json:"H_date"`
		Pilot        string    `json:"pilot"`
		Glider       string    `json:"glider"`
		Glider_id    string    `json:"glider_id"`
		Track_length float64   `json:"track_length"`
	*/

	// TESTING
	testTrack := trackSlice[1]
	switch field {
	case "H_date":
		//w.Write([]byte(testTrack.H_date))
	case "pilot":
		w.Write([]byte(testTrack.Pilot))
	case "glider":
		w.Write([]byte(testTrack.Glider))
	case "glider_id":
		w.Write([]byte(testTrack.Glider_id))
	case "track_length":
		//w.Write([]byte(testTrack.Track_length))
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

// ***************************************************
// Main:

func main() {

	// Uses mux for regex.
	router := mux.NewRouter()

	// Functions to handle the URL paths.
	router.HandleFunc("/igcinfo/api", getApiInfo)                                          // Done
	router.HandleFunc("/igcinfo/api/igc", handleTracks)                                    // Done
	router.HandleFunc("/igcinfo/api/igc/{id:[0-9]+}", getTrackByID)                        // Done
	router.HandleFunc("/igcinfo/api/igc/{id:[0-9]+}/{field:[a-z-A-Z]+}", getDetailedTrack) //

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
