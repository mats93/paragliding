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

/* Replaced by mongodb method.
// Function to test: retriveTrackByID().
// Test if the correct ID is returned.
func Test_retriveTrackByID(t *testing.T) {

	// Check if it returns an error if the trackSlice is emtpy.
	if _, err := retriveTrackByID(10); err == nil {
		t.Error("Function did not return error when track array was empty")
	}

	// Adds a track to the trackSlice with id 10.
	trackSlice = append(trackSlice,
		track{10, time.Now(), "pilot", "glider", "glider_id", 20.4, "http://test.test"})

	// Check if it returns an error if the track was not found.
	if _, err := retriveTrackByID(1); err == nil {
		t.Error("Function did not return error when requestet ID did not exist")
	}

	// Check if the correct track with the specified ID is sent back.
	newTrack, _ := retriveTrackByID(10)
	if newTrack.ID != 10 {
		t.Errorf("Function did not return the correct track: got %d want %d",
			newTrack.ID, 10)
	}

	// Sets the trackSlice to nil (empty).
	trackSlice = nil
} */

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
// Test to check the returned status code, content-type and data for the function.
func Test_HandleTracks(t *testing.T) {
	// Connects the the database.
	database := mongodb.DatabaseInit(COLLECTION)

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
		// Check the status code is what we expect (200).
		status := recorder.Code
		if status != http.StatusNotFound {
			t.Errorf("Handler returned wrong status code: got %v want %v",
				status, http.StatusNotFound)
		}
		// Expected resulsts when the DB is empty.
		expectedWhenEmpty := "[]"

		// Checks if data returned is the same as expected (empty array).
		if actual != expectedWhenEmpty {
			t.Errorf("Handler returned wrong data: got %v want %v",
				actual, expectedWhenEmpty)
		}
	} else {
		// Check the status code is what we expect (200).
		status := recorder.Code
		if status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
		// Slice of ints, to hold the IDs.
		var intSlice []int

		// Gets all tracks from the database.
		// Loops through them, appending their ID to the new slice.
		tracks, _ := database.FindAll()
		for i := 0; i < len(tracks); i++ {
			intSlice = append(intSlice, tracks[i].ID)
		}
		// Converts the struct to json.
		expected, _ := json.Marshal(intSlice)

		// Checks if data returned is the same as the contents of the DB.
		if actual != string(expected) {
			t.Errorf("Handler returned wrong data: got %v want %v",
				actual, string(expected))
		}
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

	// Check if the track was added to the slice.
	if trackSlice != nil && trackSlice[0].ID != 1 {
		t.Error("The track was not added in the trackSlice")
	}

	// Sets the trackSlice to nil (empty).
	trackSlice = nil
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

	// Adding test data to compare with, and adds it to the slice.
	trackTest := mongodb.Track{1, time.Now(), "pilot1", "glider1", "glider_id1", 21, "http://test.test"}
	trackSlice = append(trackSlice, trackTest)

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

	if actual != string(expected) {
		t.Errorf("Handler returned wrong data: got \"%v\" want \"%v\"",
			actual, string(expected))
	}

	// Sets the trackSlice to nil (empty).
	trackSlice = nil
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

	// Adding test data.
	trackSlice = append(trackSlice, mongodb.Track{1, time.Now(), "pilot1", "glider1", "glider_id1", 21, "http://test.test"})

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

	// Sets the trackSlice to nil (empty).
	trackSlice = nil
}

// Function to test: GetDetailedTrack().
// Test to check the returned status code, content-type and the data for the function.
func Test_GetDetailedTrack(t *testing.T) {

	// Adding test data to compare with, and adds it to the slice.
	expectedPilot := "pilot1"
	trackSlice = append(trackSlice, mongodb.Track{1, time.Now(), expectedPilot, "glider1", "glider_id1", 21, "http://test.test"})

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

	// Sets the trackSlice to nil (empty).
	trackSlice = nil
}
