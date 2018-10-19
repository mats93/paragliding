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

	"github.com/gorilla/mux"
	"github.com/mats93/paragliding/mongodb"
)

// Cant test DeleteAllTracks, since it wil delete all tracks.

// Function to test: GetTrackCount().
// Test to check the returned status code, content-type and data for the function.
func Test_GetTrackCount(t *testing.T) {
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
	database := mongodb.DatabaseInit("Tracks")
	count, _ := database.GetCount()

	// Converts int to string.
	expected := strconv.Itoa(count)

	// Check the response body is what we expect.
	actual := recorder.Body.String()

	if actual != expected {
		t.Errorf("Handler returned wrong data: got %v want %v",
			actual, expected)
	}

	// Closes the database session.
	defer mongodb.MDB.Session.Close()
}
