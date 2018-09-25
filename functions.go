/*
	File: functions.go
  Contains functions used by the igcinfo program.
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	igc "github.com/marni/goigc"
	"github.com/rickb777/date/period"
)

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
	// ToDo -> Encode to json properly.
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
		// Sets header status code to 400 "Bad request", and returns error message.
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error: Malformed POST request, should be '{\"url\": \"<url>\"}'"))

		// If the json decoding works:
	} else {
		// Adds the new track from the URL.
		trackFile, err := igc.ParseLocation(newUrl.Url)

		// If the igc parser fails:
		if err != nil {
			// Sets header status code to 400 "Bad request", and returns error message.
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error: Bad url, could not parse the IGC data"))

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
			// ToDo -> Encode to json properly.
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
