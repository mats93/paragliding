/*
	File: track.go
  Contains functions used by API calls to the "Track paths".
*/

package track

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	igc "github.com/marni/goigc"
	"github.com/mats93/paragliding/mongodb"
	"github.com/mats93/paragliding/webhook"
	"github.com/rickb777/date/period"
)

// INFORMATION = The API inforomation text.
const INFORMATION = "Service for Paragliding tracks"

// VERSION = The API version that is currently running.
const VERSION = "v1"

// Format for the API information.
type apiInfo struct {
	Uptime  string `json:"uptime"`
	Info    string `json:"info"`
	Version string `json:"version"`
}

// Format for the url information.
type url struct {
	URL string `json:"url"`
}

// Format for the id to be returned when a new track is posted.
type id struct {
	ID int `json:"id"`
}

// The start time for the API service. Gets injected from main.
var StartTime time.Time

// The MongoDB collection to use. Gets injected from main or test.
var Collection string

// Redirects to the /paragliding/api.
func RedirectToInfo(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, r.RequestURI+"/api", 301)
}

// GET: Returns information about the API.
// Output: application/json
func GetAPIInfo(w http.ResponseWriter, r *http.Request) {
	// Calculates the duration since application start.
	// Uses the ISO 8601 Duration format.
	// The Date package "github.com/rickb777/date/period" is used for this.
	duration := time.Since(StartTime)
	p, _ := period.NewOf(duration)
	timeStr := p.String()

	// Creates a new struct for the API info.
	currentAPI := apiInfo{timeStr, INFORMATION, VERSION}

	// Converts the struct to json.
	json, err := json.Marshal(currentAPI)
	if err != nil {
		// Sets header status code to 500 "Internal server error" and logs the error.
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
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
	// Connects to the database.
	database := mongodb.DatabaseInit(Collection)

	// Gets the Count of Tracks in the DB.
	count, _ := database.GetCount()

	// Check if there are any tracks in the DB.
	if count != 0 {
		// Slice of ints, to hold the IDs.
		var idSlice []int

		// Gets all tracks from the database.
		// Loops through them, appending their ID to the new slice.
		tracks, _ := database.FindAll()
		for i := 0; i < len(tracks); i++ {
			idSlice = append(idSlice, tracks[i].ID)
		}

		// Converts the struct to json.
		json, err := json.Marshal(idSlice)
		if err != nil {
			// Sets header status code to 500 "Internal server error" and logs the error.
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
		} else {
			// Sets header content-type to application/json and status code to 200 (OK).
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			// Returns the array of IDs.
			w.Write([]byte(json))
		}
	} else {
		// There are no tracks stored in the DB.
		// Sets header content-type to application/json and status code to 404 (Not found).
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		// Returns an empty array.
		w.Write([]byte("[]"))
	}
}

// POST: Takes the post request as json format and inserts a new track to the DB.
// Input/Output: application/json
func insertNewTrack(w http.ResponseWriter, r *http.Request) {
	// Sets the collection to be used in the webhook package.
	webhook.CollectionTrack = Collection
	webhook.CollectionWebhook = "Webhooks"
	// Connects to the database.
	database := mongodb.DatabaseInit(Collection)

	var newURL url

	// Decodes the json url and converts it to a struct.
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newURL)

	if err != nil {
		// The decoding failed.
		// Sets header status code to 400 "Bad request", and returns error message.
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error: Malformed POST request, should be '{\"url\": \"<url>\"}'"))

	} else {
		// The decoding was sucessful.
		// Adds the new track provided by the POST request.
		trackFile, err := igc.ParseLocation(newURL.URL)

		if err != nil {
			// The igc parser failed.
			// Sets header status code to 400 "Bad request", and returns error message.
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error: Bad url, could not parse the IGC data"))

		} else {
			// The igc parser worked.
			// Calculates the total distance for the track.
			var sum float64
			// Loops through all Points[] in the track.
			for i := 0; i < len(trackFile.Points)-1; i++ {
				// Adds the distance between two points togheter.
				sum += trackFile.Points[i].Distance(trackFile.Points[i+1])
			}

			// This part is in critical sector, should be locked by mutex if threads are used.
			// Generates a new ID to be used.
			newID := database.GetNewID()

			// Generates a timestamp for the track.
			timeStamp := mongodb.GenerateTimestamp()

			// Adds the new track to the database.
			database.Insert(mongodb.Track{newID, timeStamp, trackFile.Header.Date, trackFile.Pilot, trackFile.GliderType,
				trackFile.GliderID, sum, newURL.URL})

			// Critical sector ends.
			// Check if any webhooks needs to be notified of changes.
			webhook.CheckWebhooks()

			// Converts the id to json format by using the id struct and Marshaling the struct to json.
			idStruct := id{newID}
			json, err := json.Marshal(idStruct)
			if err != nil {
				// Sets header status code to 500 "Internal server error" and logs the error.
				w.WriteHeader(http.StatusInternalServerError)
				log.Fatal(err)
			} else {
				// Sets header content-type to application/json and status code to 200 (OK).
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)

				// Returns the given tracks ID as json.
				w.Write([]byte(json))
			}
		}
	}
}

// POST, GET: Track registration.
// Input/Output: application/json
func HandleTracks(w http.ResponseWriter, r *http.Request) {
	// Calls functions to handle the GET and POST requests.
	switch r.Method {
	case "GET":
		allTrackIDs(w, r)

	case "POST":
		insertNewTrack(w, r)
	}
}

// GET: Returns metadata about a given track with the provided '<id>'.
// Output: application/json
func GetTrackByID(w http.ResponseWriter, r *http.Request) {
	// Connects to the database.
	database := mongodb.DatabaseInit(Collection)

	var id int
	// Gets the ID from the URL and converts it to an integer.
	fmt.Sscanf(r.URL.Path, "/paragliding/api/track/%d", &id)

	// Tries to retrive the track with the requestet ID.
	// Retrieves the ID from the url.
	rTrackSlice, err := database.FindByID(id)
	if err == nil {
		// The request is valid, the track was found.
		rTrack := rTrackSlice[0]

		// Converts the struct to json and outputs it.
		json, err := json.Marshal(rTrack)
		if err != nil {
			// Sets header status code to 500 "Internal server error" and logs the error.
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
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
func GetDetailedTrack(w http.ResponseWriter, r *http.Request) {
	// Connects to the database.
	database := mongodb.DatabaseInit(Collection)

	var id int
	var field string

	// Gets the field and ID from the URL, converts the ID to an integer.
	fmt.Sscanf(r.URL.Path, "/paragliding/api/track/%d/%s", &id, &field)

	// Tries to retrive the track with the requestet ID.
	// Retrieves the ID from the url.
	rTrack, err := database.FindByID(id)
	if err == nil {
		// The request is valid, the track was found.

		// Retrieves the field specified, or 404 field not found.
		output := ""
		switch field {
		case "H_date":
			// Converts time.Time to string with string() method.
			output = rTrack[0].HDate.String()
		case "pilot":
			output = rTrack[0].Pilot
		case "glider":
			output = rTrack[0].Glider
		case "glider_id":
			output = rTrack[0].GliderID
		case "track_length":
			// Converts float64 to string.
			output = strconv.FormatFloat(rTrack[0].TrackLength, 'f', 6, 64)
		case "track_src_url":
			output = rTrack[0].TrackSrcURL
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
