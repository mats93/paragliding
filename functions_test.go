/*
  File: functions_test.go
  Contains unit tests for functions.go
*/

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Test if the correct ID is returned.
func Test_retriveTrackById(t *testing.T) {

	// Check if it returns an error if the trackSlice is emtpy.
	if _, err := retriveTrackById(10); err == nil {
		t.Error("Function did not return error when track array was empty")
	}

	// Adds a track to the trackSlice.
	trackSlice = append(trackSlice,
		track{10, time.Now(), "pilot", "glider", "glider_id", 20.4})

	// Check if it returns an error if the track was not found.
	if _, err := retriveTrackById(1); err == nil {
		t.Error("Function did not return error when requestet ID did not excist")
	}

	// Check if the correct track with the specified ID is sent back.
	newTrack, _ := retriveTrackById(10)
	if newTrack.Id != 10 {
		t.Errorf("Function did not return the correct track: got %d want %d",
			newTrack.Id, 10)
	}

	// Cleaning up.
	trackSlice = nil
}

// Test to check the returned status code, content-type and data for the function.
// Tests the GET request with zero tracks in memory.
func Test_handleTracks_GET_NoTracks(t *testing.T) {

	// Starts a http test server to record the output for the function we want to test.
	req, err := http.NewRequest("GET", "/igcinfo/api/igc", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Starts a http test server to record the output for the function we want to test.
	// This function is located in the 'handler.go' file.
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(handleTracks)
	handler.ServeHTTP(recorder, req)

	// Check the status code is what we expect (200).
	status := recorder.Code
	if status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check if the content-type is what we expect (application/json).
	content := recorder.HeaderMap.Get("content-type")
	if content != "application/json" {
		t.Errorf("Handler returned wrong content-type: got %s want %s",
			content, "application/json")
	}

	// Check the response body is what we expect, when 0 elements are in memory.
	expected := "[]"
	actual := recorder.Body.String()

	if actual != expected {
		t.Errorf("Handler returned wrong data: got %v want %v",
			actual, expected)
	}
}

// Tests the GET request with 3 tracks in memory.
func Test_handleTracks_GET_WithTracks(t *testing.T) {

	// Adding test data to compare with
	trackSlice = append(trackSlice,
		track{1, time.Now(), "pilot1", "glider1", "glider_id1", 21},
		track{2, time.Now(), "pilot2", "glider2", "glider_id2", 22},
		track{3, time.Now(), "pilot3", "glider3", "glider_id3", 23})

	// Starts a http test server to record the output for the function we want to test.
	req, err := http.NewRequest("GET", "/igcinfo/api/igc", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Starts a http test server to record the output for the function we want to test.
	// This function is located in the 'handler.go' file.
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(handleTracks)
	handler.ServeHTTP(recorder, req)

	// Check the response body is what we expect, when 0 elements are in memory.
	expected := "[1, 2, 3]"
	actual := recorder.Body.String()

	if actual != expected {
		t.Errorf("Handler returned wrong data: got %v want %v",
			actual, expected)
	}
}

// Test to check the returned status code, content-type and data for the function.
// Tests the POST request to insert a new track.
func Test_handleTracks_POST(t *testing.T) {
	// Send ULR in wrong format, check for correct error.
	// Send URL in correct format, but make parser fail, check for correct error.

	// Send a correct POST request, check if data is added correctly in trackSlice.
	// Check if the respond json ID is correctly formated.

	// Check response status code and content-type.
}
