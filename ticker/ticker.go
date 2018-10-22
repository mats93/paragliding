/*
	File: ticker.go
  Contains functions used by API calls to the "Ticker paths".
*/

package ticker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/mats93/paragliding/mongodb"
)

// Emulates paging of 5.
const CAP = 5

// Timestamp struct.
type Timestamp struct {
	TimeLatest int64         `json:"t_latest"`
	TimeStart  int64         `json:"t_start"`
	TimeStop   int64         `json:"t_stop"`
	Tracks     []int         `json:"tracks"`
	Processing time.Duration `json:"processing"`
}

// Collection is the MongoDB collection to use. Gets injected from main or test.
var Collection string

// GetLastTimestamp - GET: Returns the timestamp of the last added track.
// Output: text/plain
func GetLastTimestamp(w http.ResponseWriter, r *http.Request) {
	// Connects to the database.
	database := mongodb.DatabaseInit(Collection)

	// Gets all tracks from the DB.
	tracks, err := database.FindAll()
	if err != nil || tracks == nil {
		// Sets header content-type to text/plain and status code to 204 (No content).
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNoContent)

		// Outputs error.
		w.Write([]byte("There are no tracks in the database"))
	} else {
		// Sets header content-type to text/plain and status code to 200 (OK).
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)

		// Sorts the track from high to low.
		sortedTracks := mongodb.SortTrackByTimestamp(tracks)

		// Converts int64 to string.
		timestamp := strconv.FormatInt(sortedTracks[0].Timestamp, 10)

		// Returns the last added timestamp.
		w.Write([]byte(timestamp))
	}
}

// GetTimestamps - GET: Returns timestamp information.
// Output: application/json
func GetTimestamps(w http.ResponseWriter, r *http.Request) {
	// Start time of the request.
	start := time.Now()

	// Connects to the database.
	database := mongodb.DatabaseInit(Collection)

	// Gets all tracks from the DB.
	tracks, err := database.FindAll()
	if err != nil || tracks == nil {
		// Sets header content-type to application/json and status code to 204 (No content).
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)

		// Outputs error.
		w.Write([]byte("There are no tracks in the database"))
	} else {
		var newTimestamp Timestamp

		// Lengt of the track slice (-1 = last index used).
		len := len(tracks) - 1

		// Last added timestamp (Last index in the DB).
		newTimestamp.TimeLatest = tracks[len].Timestamp

		// Sorts the tracks by timestamps, last added in index 0, first added in last index.
		sortedTracks := mongodb.SortTrackByTimestamp(tracks)

		// Adds timestamps to Timestamp struct.
		newTimestamp.TimeStart = sortedTracks[0].Timestamp
		newTimestamp.TimeStop = sortedTracks[len].Timestamp

		// The max allowed IDs to be returned.
		var maxLoops int
		if len+1 > CAP {
			maxLoops = CAP - 1
		} else {
			maxLoops = len
		}

		// Adds the IDs to the slice.
		for i := 0; i <= maxLoops; i++ {
			newTimestamp.Tracks = append(newTimestamp.Tracks, tracks[i].ID)
		}

		// Adds the processing time.
		newTimestamp.Processing = time.Since(start)

		// Converts the struct to json.
		json, err := json.Marshal(newTimestamp)
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
}

// GetTimestampsNewerThen - GET: Returns timestamp information of newer timestamps then provided.
// Output: application/json
func GetTimestampsNewerThen(w http.ResponseWriter, r *http.Request) {
	// Start time of the request.
	start := time.Now()

	var ts int64

	// Gets the timestamp from the URL.
	fmt.Sscanf(r.URL.Path, "/paragliding/api/ticker/%d", &ts)

	// Connects to the database.
	database := mongodb.DatabaseInit(Collection)

	// Gets all tracks from the DB.
	tracks, err := database.FindTrackHigherThen(ts)
	if err != nil || tracks == nil {
		// Sets header content-type to application/json and status code to 204 (No content).
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)

		// Outputs error.
		w.Write([]byte("There are no tracks in the database"))
	} else {
		var newTimestamp Timestamp

		// Lengt of the track slice (-1 = last index used).
		len := len(tracks) - 1

		// Last added timestamp (Last index in the DB).
		newTimestamp.TimeLatest = tracks[len].Timestamp

		// Sorts the tracks by timestamps, last added in index 0, first added in last index.
		sortedTracks := mongodb.SortTrackByTimestamp(tracks)

		// Adds timestamps to Timestamp struct.
		newTimestamp.TimeStart = sortedTracks[0].Timestamp
		newTimestamp.TimeStop = sortedTracks[len].Timestamp

		// The max allowed IDs to be returned.
		var maxLoops int
		if len+1 > CAP {
			maxLoops = CAP - 1
		} else {
			maxLoops = len
		}

		// Adds the IDs to the slice.
		for i := 0; i <= maxLoops; i++ {
			newTimestamp.Tracks = append(newTimestamp.Tracks, tracks[i].ID)
		}

		// Adds the processing time.
		newTimestamp.Processing = time.Since(start)

		// Converts the struct to json.
		json, err := json.Marshal(newTimestamp)
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
}
