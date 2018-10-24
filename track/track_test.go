/*
  File: track_test.go
  Contains unit tests for track.go
*/

package track

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/mats93/paragliding/mongodb"
	"github.com/rickb777/date/period"
)

// Function to test: GetAPIInfo().
// Test to check the returned status code, content-type and data for the function.
func Test_GetAPIInfo(t *testing.T) {

	// Creates a request that is passed to the handler.
	request, _ := http.NewRequest("GET", "/paragliding/api", nil)

	// Creates the recorder and router.
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()

	// Tests the function.
	router.HandleFunc("/paragliding/api", GetAPIInfo).Methods("GET")
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

	// Calculates the start time of the app.
	duration := time.Since(StartTime)
	p, _ := period.NewOf(duration)
	timeStr := p.String()

	// Excpected result from the API call.
	expected := apiInfo{timeStr, INFORMATION, VERSION}

	// Adding test data to compare with.
	var actual apiInfo
	decoder := json.NewDecoder(recorder.Body)
	decoder.Decode(&actual)

	// Check the response body is what we expect.
	if actual != expected {
		t.Errorf("Handler returned wrong data: got %v want %v",
			actual, expected)
	}
}

// Function to test: HandleTracks().
// Test to check the returned status code, content-type and data for the function when the DB is empty.
func Test_HandleTracks_EmptyDB(t *testing.T) {
	// Injects the MongoDB collection to use.
	Collection = "Tests"

	// Connects the the database.
	database := mongodb.DatabaseInit(Collection)

	// Creates a request that is passed to the handler.
	request, _ := http.NewRequest("GET", "/paragliding/api/track", nil)

	// Creates the recorder and router.
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()

	// Tests the function.
	router.HandleFunc("/paragliding/api/track", HandleTracks).Methods("GET")
	router.ServeHTTP(recorder, request)

	// Check if the content-type is what we expect (application/json).
	content := recorder.HeaderMap.Get("content-type")
	if content != "application/json" {
		t.Errorf("Handler returned wrong content-type: got %s want %s",
			content, "application/json")
	}

	// The actual retuned data.
	actual := recorder.Body.String()

	// Gets the Count of Tracks in the DB.
	count, _ := database.GetCount()

	// Check if there are any tracks in the DB.
	if count == 0 {
		// Check the status code is what we expect (400).
		status := recorder.Code
		if status != http.StatusNotFound {
			t.Errorf("Handler returned wrong status code: got %v want %v",
				status, http.StatusNotFound)
		}

		// Expected resulsts when the DB is empty.
		expected := "[]"
		// Checks if data returned is the same as expected (empty array).
		if actual != expected {
			t.Errorf("Handler returned wrong data: got %v want %v",
				actual, expected)
		}
	} else {
		t.Error("Database count is not 0, when 0 tracks are in the DB")
	}
}

// Function to test: HandleTracks().
// Test to check the returned status code, content-type and data for the function.
func Test_HandleTracks(t *testing.T) {
	// Injects the MongoDB collection to use.
	Collection = "Tests"

	// Connects the the database and inserts 3 tracks.
	database := mongodb.DatabaseInit(Collection)
	database.Insert(mongodb.Track{1, 11, time.Now(), "pilot1", "glider1", "glider_id1", 20.1, "http://test1.test"})
	database.Insert(mongodb.Track{2, 12, time.Now(), "pilot2", "glider2", "glider_id2", 20.2, "http://test2.test"})
	database.Insert(mongodb.Track{3, 13, time.Now(), "pilot3", "glider3", "glider_id3", 20.3, "http://test3.test"})

	// Expected return when 3 tracks are in the DB.
	expected := "[1,2,3]"

	// Creates a request that is passed to the handler.
	request, _ := http.NewRequest("GET", "/paragliding/api/track", nil)

	// Creates the recorder and router.
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()

	// Tests the function.
	router.HandleFunc("/paragliding/api/track", HandleTracks).Methods("GET")
	router.ServeHTTP(recorder, request)

	// Check if the content-type is what we expect (application/json).
	content := recorder.HeaderMap.Get("content-type")
	if content != "application/json" {
		t.Errorf("Handler returned wrong content-type: got %s want %s",
			content, "application/json")
	}

	// The actual retuned data.
	actual := recorder.Body.String()

	// Gets the Count of Tracks in the DB.
	count, _ := database.GetCount()

	// Check if there are any tracks in the DB.
	if count == 0 {
		t.Error("Database count is 0, when 3 tracks are insertet into the DB")
	} else {
		// Check the status code is what we expect (200).
		status := recorder.Code
		if status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		// Checks if data returned is the same as the contents of the DB.
		if actual != expected {
			t.Errorf("Handler returned wrong data: got %v want %v",
				actual, expected)
		}

		// Removes the test data.
		database.DeleteAll()
	}
}

// Function to test: HandleTracks().
// Test to check the returned status code, content-type and data when the POST request has wrong json format.
func Test_HandleTracks_POST_MalformedPost(t *testing.T) {

	// Creates a malformed (wrong json format) POST request that is passed to the handler.
	request, _ := http.NewRequest("POST", "/paragliding/api/track", strings.NewReader("wrong"))

	// Creates the recorder and router.
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()

	// Tests the function.
	router.HandleFunc("/paragliding/api/track", HandleTracks).Methods("POST")
	router.ServeHTTP(recorder, request)

	// Check the status code is what we expect (400).
	status := recorder.Code
	if status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	// Check the response body is what we expect, an error.
	expected := "Error: Malformed POST request, should be '{\"url\": \"<url>\"}'"
	actual := recorder.Body.String()

	if actual != expected {
		t.Errorf("Handler returned wrong data: got \"%v\" want \"%v\"",
			actual, expected)
	}
}

// Function to test: HandleTracks().
// Test to check the returned status code, content-type and data when the POST request has wrong url, but correct json format.
func Test_HandleTracks_POST_WrongFile(t *testing.T) {

	// POST data, correct json format, but it's not an igc file.
	postString := "{\"url\":\"wrong\"}"

	// Creates a POST request with an url that does not work that is passed to the handler.
	request, _ := http.NewRequest("POST", "/paragliding/api/track", strings.NewReader(postString))

	// Creates the recorder and router.
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()

	// Tests the function.
	router.HandleFunc("/paragliding/api/track", HandleTracks).Methods("POST")
	router.ServeHTTP(recorder, request)

	// Check the status code is what we expect (400).
	status := recorder.Code
	if status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	// Check the response body is what we expect, an error.
	expected := "Error: Bad url, could not parse the IGC data"
	actual := recorder.Body.String()

	if actual != expected {
		t.Errorf("Handler returned wrong data: got \"%v\" want \"%v\"",
			actual, expected)
	}
}

// Function to test: HandleTracks().
// Test to check the returned status code, content-type and data for the function.
func Test_HandleTracks_POST(t *testing.T) {
	// Injects the MongoDB collection to use.
	Collection = "Tests"

	// Connects the the database.
	database := mongodb.DatabaseInit(Collection)

	// Gets the current count of the DB. Should be 0.
	count, _ := database.GetCount()

	// POST data.
	postString := "{\"url\":\"http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc\"}"

	// Creates a POST request that is passed to the handler.
	request, _ := http.NewRequest("POST", "/paragliding/api/track", strings.NewReader(postString))

	// Creates the recorder and router.
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()

	// Tests the function.
	router.HandleFunc("/paragliding/api/track", HandleTracks).Methods("POST")
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

	// Check the response body is what we expect.
	expectedReturn := "{\"id\":1}"
	actualReturn := recorder.Body.String()

	if actualReturn != expectedReturn {
		t.Errorf("Handler returned wrong data: got %v want %v",
			actualReturn, expectedReturn)
	}

	// Check if the track was added to the database, newCount should be 1.
	newCount, _ := database.GetCount()
	if count+1 != newCount {
		t.Errorf("Database count is wrong: got %d want %d",
			newCount, count+1)
	}

	// Removes the test data.
	database.DeleteAll()
}

// Function to test: GetTrackByID().
// Test to check the returned status code, content-type and data when the requested track does not exist.
func Test_GetTrackByID_NoTrackExists(t *testing.T) {
	// Creates a request that is passed to the handler.
	request, _ := http.NewRequest("GET", "/paragliding/api/track/1", nil)

	// Creates the recorder and router.
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()

	// Tests the function.
	router.HandleFunc("/paragliding/api/track/1", GetTrackByID).Methods("GET")
	router.ServeHTTP(recorder, request)

	// Check the status code is what we expect (404).
	status := recorder.Code
	if status != http.StatusNotFound {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

// Function to test: GetTrackByID().
// Test to check the returned status code, content-type and data for the function.
func Test_GetTrackByID(t *testing.T) {
	// Injects the MongoDB collection to use.
	Collection = "Tests"

	// Connects the the database, and adds test data to the DB.
	database := mongodb.DatabaseInit(Collection)
	trackTest := mongodb.Track{1, 11, time.Now(), "pilot1", "glider1", "glider_id1", 21, "http://test.test"}
	database.Insert(trackTest)

	// Creates a request that is passed to the handler.
	request, _ := http.NewRequest("GET", "/paragliding/api/track/1", nil)

	// Creates the recorder and router.
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()

	// Tests the function.
	router.HandleFunc("/paragliding/api/track/1", GetTrackByID).Methods("GET")
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

	// Check the response body is what we expect.
	expected, _ := json.Marshal(trackTest)
	actual := recorder.Body.String()

	// Because of time.Time the comparrsing fails when comparing the json objects.
	// Workaround, split it and compare the pilot in a json format ("pilot":"name").
	expectedSplit := strings.Split(string(expected), ",")
	actualSplit := strings.Split(actual, ",")
	actualPilotJson := actualSplit[1]
	expectedPilotJson := expectedSplit[1]

	if actualPilotJson != expectedPilotJson {
		t.Errorf("Handler returned wrong data: got %s want %s",
			actualPilotJson, expectedPilotJson)
	}

	// Removes the test data.
	database.DeleteAll()
}

// Function to test: GetDetailedTrack().
// Test to check the returned status code, content-type when the requested ID does not exist.
func Test_GetDetailedTrack_WrongID(t *testing.T) {

	// Creates a request that is passed to the handler.
	request, _ := http.NewRequest("GET", "/paragliding/api/track/10/pilot", nil)

	// Creates the recorder and router.
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()

	// Tests the function.
	router.HandleFunc("/paragliding/api/track/10/pilot", GetDetailedTrack).Methods("GET")
	router.ServeHTTP(recorder, request)

	// Check the status code is what we expect (404).
	status := recorder.Code
	if status != http.StatusNotFound {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

// Function to test: GetDetailedTrack().
// Test to check the returned status code, content-type when a non existent field is passed.
func Test_GetDetailedTrack_WrongField(t *testing.T) {
	// Injects the MongoDB collection to use.
	Collection = "Tests"

	// Connects the the database, and adds test data to the DB.
	database := mongodb.DatabaseInit(Collection)
	database.Insert(mongodb.Track{1, 11, time.Now(), "pilot1", "glider1", "glider_id1", 21, "http://test.test"})

	// Creates a request that is passed to the handler.
	request, _ := http.NewRequest("GET", "/paragliding/api/track/1/feil", nil)

	// Creates the recorder and router.
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()

	// Tests the function.
	router.HandleFunc("/paragliding/api/track/1/feil", GetDetailedTrack).Methods("GET")
	router.ServeHTTP(recorder, request)

	// Check the status code is what we expect (200).
	status := recorder.Code
	if status != http.StatusNotFound {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

	// Removes the test data.
	database.DeleteAll()
}

// Function to test: GetDetailedTrack().
// Test to check the returned status code, content-type and the data for the function.
func Test_GetDetailedTrack(t *testing.T) {
	// Injects the MongoDB collection to use.
	Collection = "Tests"

	// Expected pilot to be returned.
	expectedPilot := "pilot1"

	// Connects the the database, and adds test data to the DB.
	database := mongodb.DatabaseInit(Collection)
	database.Insert(mongodb.Track{1, 11, time.Now(), expectedPilot, "glider1", "glider_id1", 21, "http://test.test"})

	// Creates a request that is passed to the handler.
	request, _ := http.NewRequest("GET", "/paragliding/api/track/1/pilot", nil)

	// Creates the recorder and router.
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()

	// Tests the function.
	router.HandleFunc("/paragliding/api/track/1/pilot", GetDetailedTrack).Methods("GET")
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

	// Check the response body is what we expect.
	actual := recorder.Body.String()

	if actual != expectedPilot {
		t.Errorf("Handler returned wrong data: got %v want %v",
			actual, expectedPilot)
	}

	// Removes the test data.
	database.DeleteAll()
}
