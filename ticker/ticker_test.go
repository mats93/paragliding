/*
  File: ticker_test.go
  Contains unit tests for ticker.go
*/

package ticker

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/mats93/paragliding/mongodb"
)

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

// Function to test: GetLastTimestamp().
// Test to check the returned status code, content-type and data for the function, when the DB is empty.
func Test_GetTimestamps_Empty(t *testing.T) {
	// Injects the MongoDB collection to use.
	Collection = "Tests"

	// Creates a request that is passed to the handler.
	request, _ := http.NewRequest("GET", "/paragliding/api/ticker/", nil)

	// Creates the recorder and router.
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()

	// Tests the function.
	router.HandleFunc("/paragliding/api/ticker/", GetTimestamps).Methods("GET")
	router.ServeHTTP(recorder, request)

	// Check the status code is what we expect (204).
	status := recorder.Code
	if status != http.StatusNoContent {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	// Check if the content-type is what we expect (application/json).
	content := recorder.HeaderMap.Get("content-type")
	if content != "application/json" {
		t.Errorf("Handler returned wrong content-type: got %s want %s",
			content, "application/json")
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
func Test_GetTimestamps(t *testing.T) {
	// Injects the MongoDB collection to use.
	Collection = "Tests"

	// Connects the the database and inserts 3 tracks.
	database := mongodb.DatabaseInit(Collection)
	database.Insert(mongodb.Track{1, 111, time.Now(), "pilot1", "glider1", "glider_id1", 20.1, "http://test1.test"})
	database.Insert(mongodb.Track{2, 222, time.Now(), "pilot2", "glider2", "glider_id2", 20.2, "http://test2.test"})
	database.Insert(mongodb.Track{3, 333, time.Now(), "pilot3", "glider3", "glider_id3", 20.3, "http://test3.test"})

	// Creates a request that is passed to the handler.
	request, _ := http.NewRequest("GET", "/paragliding/api/ticker/", nil)

	// Creates the recorder and router.
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()

	// Tests the function.
	router.HandleFunc("/paragliding/api/ticker/", GetTimestamps).Methods("GET")
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

	// Expected timestamp to be returned.
	expectedLatestJson := "{\"t_latest\":333"
	expectedStartJson := "\"t_start\":333"
	expectedStopJson := "\"t_stop\":111"
	expectedIDsJson := "\"tracks\":[1"

	// The actual retuned data.
	// Dont want to compare against processing. Splits the json and compers with it.
	actual := recorder.Body.String()
	actualSplit := strings.Split(actual, ",")
	actualLatestJson := actualSplit[0]
	actualStartJson := actualSplit[1]
	actualStopJson := actualSplit[2]
	actualIDsJson := actualSplit[3]

	if actualLatestJson != expectedLatestJson {
		t.Errorf("Function returned wrong Latest field: got %s want %s",
			actualLatestJson, expectedLatestJson)
	}
	if actualStartJson != expectedStartJson {
		t.Errorf("Function returned wrong Start field: got %s want %s",
			actualStartJson, expectedStartJson)
	}
	if actualStopJson != expectedStopJson {
		t.Errorf("Function returned wrong Stop field: got %s want %s",
			actualStopJson, expectedStopJson)
	}
	if actualIDsJson != expectedIDsJson {
		t.Errorf("Function returned wrong IDs field: got %s want %s",
			actualIDsJson, expectedIDsJson)
	}

	// Removes the test data.
	database.DeleteAllTracks()
}

// Function to test: GetTimestampsNewerThen().
// Test to check the returned status code, content-type and data for the function, when the DB is empty.
func Test_GetTimestampsNewerThen_Empty(t *testing.T) {
	// Injects the MongoDB collection to use.
	Collection = "Tests"

	// Creates a request that is passed to the handler.
	request, _ := http.NewRequest("GET", "/paragliding/api/ticker/1", nil)

	// Creates the recorder and router.
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()

	// Tests the function.
	router.HandleFunc("/paragliding/api/ticker/1", GetTimestampsNewerThen).Methods("GET")
	router.ServeHTTP(recorder, request)

	// Check the status code is what we expect (204).
	status := recorder.Code
	if status != http.StatusNoContent {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	// Check if the content-type is what we expect (application/json).
	content := recorder.HeaderMap.Get("content-type")
	if content != "application/json" {
		t.Errorf("Handler returned wrong content-type: got %s want %s",
			content, "application/json")
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

// Function to test: GetTimestampsNewerThen().
// Test to check the returned status code, content-type and data for the function.
func Test_GetTimestampsNewerThen(t *testing.T) {
	// Injects the MongoDB collection to use.
	Collection = "Tests"

	// Connects the the database and inserts 3 tracks.
	database := mongodb.DatabaseInit(Collection)
	database.Insert(mongodb.Track{1, 111, time.Now(), "pilot1", "glider1", "glider_id1", 20.1, "http://test1.test"})
	database.Insert(mongodb.Track{2, 222, time.Now(), "pilot2", "glider2", "glider_id2", 20.2, "http://test2.test"})
	database.Insert(mongodb.Track{3, 333, time.Now(), "pilot3", "glider3", "glider_id3", 20.3, "http://test3.test"})

	// Creates a request that is passed to the handler.
	request, _ := http.NewRequest("GET", "/paragliding/api/ticker/222", nil)

	// Creates the recorder and router.
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()

	// Tests the function.
	router.HandleFunc("/paragliding/api/ticker/222", GetTimestampsNewerThen).Methods("GET")
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

	// Expected timestamp to be returned.
	expectedLatestJson := "{\"t_latest\":333"
	expectedStartJson := "\"t_start\":333"
	expectedStopJson := "\"t_stop\":333"
	expectedIDsJson := "\"tracks\":[3]"

	// The actual retuned data.
	// Dont want to compare against processing. Splits the json and compers with it.
	actual := recorder.Body.String()
	actualSplit := strings.Split(actual, ",")
	actualLatestJson := actualSplit[0]
	actualStartJson := actualSplit[1]
	actualStopJson := actualSplit[2]
	actualIDsJson := actualSplit[3]

	if actualLatestJson != expectedLatestJson {
		t.Errorf("Function returned wrong Latest field: got %s want %s",
			actualLatestJson, expectedLatestJson)
	}
	if actualStartJson != expectedStartJson {
		t.Errorf("Function returned wrong Start field: got %s want %s",
			actualStartJson, expectedStartJson)
	}
	if actualStopJson != expectedStopJson {
		t.Errorf("Function returned wrong Stop field: got %s want %s",
			actualStopJson, expectedStopJson)
	}
	if actualIDsJson != expectedIDsJson {
		t.Errorf("Function returned wrong IDs field: got %s want %s",
			actualIDsJson, expectedIDsJson)
	}

	// Removes the test data.
	database.DeleteAllTracks()
}
