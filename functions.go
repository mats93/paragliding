/*
  File: functions.go
  Contains helper functions for the http handle function calls.
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	igc "github.com/marni/goigc"
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
