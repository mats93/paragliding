/*
  File: admin_test.go
  Contains unit tests for admin.go
*/

package admin

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/mats93/paragliding/mongodb"
)

// Function to test: GetTrackCount().
// Test to check the returned status code, content-type and data for the function.
func Test_GetTrackCount(t *testing.T) {
	// Injects the MongoDB collection to use.
	Collection = "TestTracks"

	// Creates a request that is passed to the handler.
	request, _ := http.NewRequest("GET", "/paragliding/admin/api/tracks_count", nil)

	// Creates the recorder and router.
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()

	// Tests the function.
	router.HandleFunc("/paragliding/admin/api/tracks_count", GetTrackCount).Methods("GET")
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

	// The count we expect to get in the body.
	database := mongodb.DatabaseInit(Collection)
	count, _ := database.GetCount()

	// Converts int to string.
	expected := strconv.Itoa(count)

	// Check the response body is what we expect.
	actual := recorder.Body.String()

	if actual != expected {
		t.Errorf("Handler returned wrong data: got %v want %v",
			actual, expected)
	}
}

// Function to test: DeleteAllTracks().
// Test to check the returned status code, content-type when the wrong method is used.
func Test_DeleteAllTracks_WrongMethod(t *testing.T) {
	// Injects the MongoDB collection to use.
	Collection = "TestTracks"

	// Creates a request that is passed to the handler.
	request, _ := http.NewRequest("GET", "/paragliding/admin/api/tracks", nil)

	// Creates the recorder and router.
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()

	// Tests the function.
	router.HandleFunc("/paragliding/admin/api/tracks", DeleteAllTracks).Methods("GET")
	router.ServeHTTP(recorder, request)

	// Check the status code is what we expect (400).
	status := recorder.Code
	if status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	// Check if the content-type is what we expect (text/plain).
	content := recorder.HeaderMap.Get("content-type")
	if content != "text/plain" {
		t.Errorf("Handler returned wrong content-type: got %s want %s",
			content, "text/plain")
	}
}

// Function to test: DeleteAllTracks().
// Test to check the returned status code, content-type and data for the function.
func Test_DeleteAllTracks(t *testing.T) {
	// Injects the MongoDB collection to use.
	Collection = "TestTracks"

	// Connets to the DB and fills it with 5 tracks.
	database := mongodb.DatabaseInit(Collection)
	database.Insert(mongodb.Track{1, 11, time.Now(), "pilot1", "glider1", "glider_id1", 20.1, "http://test1.test"})
	database.Insert(mongodb.Track{2, 12, time.Now(), "pilot2", "glider2", "glider_id2", 20.2, "http://test2.test"})
	database.Insert(mongodb.Track{3, 13, time.Now(), "pilot3", "glider3", "glider_id3", 20.3, "http://test3.test"})
	database.Insert(mongodb.Track{4, 14, time.Now(), "pilot4", "glider4", "glider_id4", 20.4, "http://test4.test"})
	database.Insert(mongodb.Track{5, 15, time.Now(), "pilot5", "glider5", "glider_id5", 20.5, "http://test5.test"})

	// Expected return for the function is 5, because 5 tracks where deletet (all).
	expected := "5"

	// Creates a request that is passed to the handler.
	request, _ := http.NewRequest("DELETE", "/paragliding/admin/api/tracks", nil)

	// Creates the recorder and router.
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()

	// Tests the function.
	router.HandleFunc("/paragliding/admin/api/tracks", DeleteAllTracks).Methods("DELETE")
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

	// Check if the fucntion returned as expected.
	actual := recorder.Body.String()
	if actual != expected {
		t.Errorf("Handler returned wrong data: got %v want %v",
			actual, expected)
	}

	// Check if the function actually deletet all tracks from the DB.
	expectedCount := 0
	actualCount, _ := database.GetCount()
	if actualCount != expectedCount {
		t.Errorf("Database returned wrong count after deletion: got %d want %d",
			actualCount, expectedCount)
	}
}
