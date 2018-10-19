/*
	File: admin.go
  Contains functions used by API calls to the "Admin paths".
*/

package admin

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/mats93/paragliding/mongodb"
)

// The database collection to use.
const COLLECTION = "Tracks"

// GET: Returns the current count of all tracks in the DB.
// Output: text/plain
func GetTrackCount(w http.ResponseWriter, r *http.Request) {
	// Connects to the database.
	database := mongodb.DatabaseInit(COLLECTION)

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
	database := mongodb.DatabaseInit(COLLECTION)

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

// ToDo: Delete this
func InsertSomeTracks(w http.ResponseWriter, r *http.Request) {
	// Connects to the database.
	database := mongodb.DatabaseInit(COLLECTION)
	// Inserts 5 tracks to the database.
	database.Insert(mongodb.Track{1, time.Now(), "pilot1", "glider1", "glider_id1", 20.1, "http://test1.test"})
	database.Insert(mongodb.Track{2, time.Now(), "pilot2", "glider2", "glider_id2", 20.2, "http://test2.test"})
	database.Insert(mongodb.Track{3, time.Now(), "pilot3", "glider3", "glider_id3", 20.3, "http://test3.test"})
	database.Insert(mongodb.Track{4, time.Now(), "pilot4", "glider4", "glider_id4", 20.4, "http://test4.test"})
	database.Insert(mongodb.Track{5, time.Now(), "pilot5", "glider5", "glider_id5", 20.5, "http://test5.test"})

	defer mongodb.MDB.Session.Close()
}
