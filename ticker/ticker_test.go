/*
  File: ticker_test.go
  Contains unit tests for ticker.go
*/

package ticker

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/mats93/paragliding/mongodb"
)

// Function to test: sortTrackByTimestamp().
// Test to check if the slice was sorted correctly.
func Test_sortTrackByTimestamp(t *testing.T) {
	// Injects the MongoDB collection to use.
	Collection = "Tests"

	// Connects the the database and inserts 3 tracks.
	// The last inserted has the highest timestamp.
	database := mongodb.DatabaseInit(Collection)
	database.Insert(mongodb.Track{1, 111, time.Now(), "pilot1", "glider1", "glider_id1", 20.1, "http://test1.test"})
	database.Insert(mongodb.Track{2, 222, time.Now(), "pilot2", "glider2", "glider_id2", 20.2, "http://test2.test"})
	database.Insert(mongodb.Track{3, 333, time.Now(), "pilot3", "glider3", "glider_id3", 20.3, "http://test3.test"})

	// Returns all tracks from the DB.
	tracks, _ := database.FindAll()

	// Try to sort the slice.
	sorted := sortTrackByTimestamp(tracks)

	// The unsorted track should have the highest timestamp in index 2.
	// The sorted track should have the highest timestamp in index 0.
	if sorted[2].Timestamp > sorted[0].Timestamp {
		t.Errorf("Function did not sort correctly")
	}

	// Removes the test data.
	database.DeleteAllTracks()
}

// Function to test: GetLastTimestamp().
// Test to check the returned status code, content-type and data for the function, when the DB is empty.
func Test_GetLastTimestamp_Empty(t *testing.T) {
	// Injects the MongoDB collection to use.
	Collection = "Tests"
	// Creates a request that is passed to the handler.
	request, _ := http.NewRequest("GET", "/paragliding/api/ticker/latest", nil)

	// Creates the recorder and router.
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()

	// Tests the function.
	router.HandleFunc("/paragliding/api/ticker/latest", GetLastTimestamp).Methods("GET")
	router.ServeHTTP(recorder, request)

	// Check the status code is what we expect (204).
	status := recorder.Code
	if status != http.StatusNoContent {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	// Check if the content-type is what we expect (text/plain).
	content := recorder.HeaderMap.Get("content-type")
	if content != "text/plain" {
		t.Errorf("Handler returned wrong content-type: got %s want %s",
			content, "text/plain")
	}

	// Expected error message.
	expected := "There are no tracks in the database"
	// The actual error message returned.
	actual := recorder.Body.String()

	if actual != expected {
		t.Errorf("Handler returned wrong data: got %s want %s",
			actual, expected)
	}
}

// Function to test: GetLastTimestamp().
// Test to check the returned status code, content-type and data for the function.
func Test_GetLastTimestamp(t *testing.T) {
	// Injects the MongoDB collection to use.
	Collection = "Tests"

	// Connects the the database and inserts 3 tracks.
	database := mongodb.DatabaseInit(Collection)
	database.Insert(mongodb.Track{1, 111, time.Now(), "pilot1", "glider1", "glider_id1", 20.1, "http://test1.test"})
	database.Insert(mongodb.Track{2, 222, time.Now(), "pilot2", "glider2", "glider_id2", 20.2, "http://test2.test"})
	database.Insert(mongodb.Track{3, 333, time.Now(), "pilot3", "glider3", "glider_id3", 20.3, "http://test3.test"})

	// Creates a request that is passed to the handler.
	request, _ := http.NewRequest("GET", "/paragliding/api/ticker/latest", nil)

	// Creates the recorder and router.
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()

	// Tests the function.
	router.HandleFunc("/paragliding/api/ticker/latest", GetLastTimestamp).Methods("GET")
	router.ServeHTTP(recorder, request)

	// Check the status code is what we expect (200).
	status := recorder.Code
	if status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check if the content-type is what we expect (text/plain).
	content := recorder.HeaderMap.Get("content-type")
	if content != "text/plain" {
		t.Errorf("Handler returned wrong content-type: got %s want %s",
			content, "text/plain")
	}

	// Expected timestamp to be returned.
	expected := "333"
	// The actual retuned data.
	actual := recorder.Body.String()

	if actual != expected {
		t.Errorf("Handler returned wrong data: got %s want %s",
			actual, expected)
	}

	// Removes the test data.
	database.DeleteAllTracks()
}
