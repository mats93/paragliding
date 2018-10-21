/*
	File: ticker.go
  Contains functions used by API calls to the "Ticker paths".
*/

package ticker

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/mats93/paragliding/mongodb"
)

// Emulates paging of 5.
const CAP = 5

// Timestamp struct.
var Timestamp struct {
	TimeLatest int64           `json:"t_latest"`
	TimeStart  int64           `json:"t_start"`
	TimeStop   int64           `json:"t_stop"`
	Tracks     map[int64]int64 `json:"tracks"`
	Processing int64           `json:"processing"`
}

// Collection is the MongoDB collection to use. Gets injected from main or test.
var Collection string

// Takes a slice of Tracks, sorts them from newest to oldest (increasing), returns the sortet slice.
func sortTrackByTimestamp(track []mongodb.Track) []mongodb.Track {
	// The function works on a buffer.
	buffer := append([]mongodb.Track(nil), track...)

	// Sorts the slice based on the Timestmap.
	sort.Slice(buffer, func(i, j int) bool {
		return buffer[i].Timestamp > buffer[j].Timestamp
	})

	// Returns the sorted track.
	return buffer
}

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
		sortedTracks := sortTrackByTimestamp(tracks)

		// Converts int64 to string.
		timestamp := strconv.FormatInt(sortedTracks[0].Timestamp, 10)

		// Returns the last added timestamp.
		w.Write([]byte(timestamp))
	}
}

// GetTimestamps - GET: Returns timestamp information.
// Output: application/json
func GetTimestamps(w http.ResponseWriter, r *http.Request) {

}

// GetTimestampsNewerThen - GET: Returns timestamp information of newer timestamps then provided.
// Output: application/json
func GetTimestampsNewerThen(w http.ResponseWriter, r *http.Request) {

}
