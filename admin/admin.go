/*
	File: admin.go
  Contains functions used by API calls to the "Admin paths".
*/

package admin

import (
	"log"
	"net/http"
	"strconv"

	"github.com/mats93/paragliding/mongodb"
)

// GET: Returns the current count of all tracks in the DB.
// Output: text/plain
func GetTrackCount(w http.ResponseWriter, r *http.Request) {
	// Connects to the database.
	database := mongodb.DatabaseInit("Tracks")

	count, err := database.GetCount()
	if err != nil {
		// Sets header status code to 500 "Internal server error" and logs the error.
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
	} else {
		// Sets header content-type to text/plain and status code to 200 (OK).
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)

		// Converts int to string.
		strCount := strconv.Itoa(count)

		// Returns the Count.
		w.Write([]byte(strCount))
	}

	// Closes the database session.
	defer mongodb.MDB.Session.Close()
}

// DELETE: Deletes all tracks in the DB.
// Output: text/plain
func DeleteAllTracks(w http.ResponseWriter, r *http.Request) {
	// Connects to the database.
	database := mongodb.DatabaseInit("Tracks")

	// Gets the current count of the database.
	count, err := database.GetCount()
	if err != nil {
		// Sets header status code to 500 "Internal server error" and logs the error.
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
	} else {
		// Deletes all tracks from the database.
		err := database.DeleteAllTracks()
		if err != nil {
			// Sets header status code to 500 "Internal server error" and logs the error.
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
		} else {
			// Sets header content-type to text/plain and status code to 200 (OK).
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)

			// Converts int to string.
			strCount := strconv.Itoa(count)

			// Returns the Count of all tracks that was removed.
			w.Write([]byte(strCount))
		}
	}

	// Closes the database session.
	defer mongodb.MDB.Session.Close()
}
