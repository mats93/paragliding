/*
  File: handler_test.go
  Contains unit tests for handler.go
*/

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

// ToDo: Refactor til å bruke mux router for all http testing.
// Test to check the returned status code, content-type and data for the function.
func Test_getApiInfo(t *testing.T) {

	// Creates a request that we pass to the handler.
	request, _ := http.NewRequest("GET", "/igcinfo/api", nil)

	// Creates the recorder and router.
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()

	// Tests the function.
	router.HandleFunc("/igcinfo/api", getApiInfo).Methods("GET")
	router.ServeHTTP(recorder, request)

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

	// Excpected result from the API call.
	expected := apiInfo{"P0D", INFORMATION, VERSION}

	// Adding test data to compare with ("P0D" is what time.Time returns).
	var actual apiInfo
	decoder := json.NewDecoder(recorder.Body)
	decoder.Decode(&actual)

	// Check the response body is what we expect.
	if actual != expected {
		t.Errorf("Handler returned wrong data: got %v want %v",
			actual, expected)
	}
}

// Test to check the returned status code, content-type and data for the function.
func Test_getTrackByID(t *testing.T) {

}

// Test to check the returned status code, content-type and data for the function.
func Test_getDetailedTrack(t *testing.T) {

}
