/*
	File: webhook_test.go
  Contains unit tests for webhook.go
*/

package webhook

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/mats93/paragliding/mongodb"
)

// Function to test: CheckWebhooks()
// Test to check if the discord channel resieves the message.
func Test_CheckWebhooks(t *testing.T) {
	CollectionTrack = "TestTracks"
	CollectionWebhook = "TestWebhooks"

	discordTestChannel := "https://discordapp.com/api/webhooks/504733605344313354/8sLrUSTCJQxB-8BAcRBH27T3jm8xWtCv1DjpjdbJnhqWqNsjxbatT_EWFsQtPnKzuuCf"

	// Adds some webhooks.
	database := mongodb.DatabaseInit(CollectionWebhook)
	database.InsertWebhook(mongodb.Webhook{discordTestChannel, 3, 0})

	databaseTracks := mongodb.DatabaseInit(CollectionTrack)

	// Add 3 tracks to the DB.
	databaseTracks.Insert(mongodb.Track{1, 111, time.Now(), "pilot1", "glider1", "glider_id1", 20.1, "http://test1.local"})
	CheckWebhooks()
	databaseTracks.Insert(mongodb.Track{2, 222, time.Now(), "pilot2", "glider2", "glider_id2", 20.2, "http://test2.local"})
	CheckWebhooks()
	databaseTracks.Insert(mongodb.Track{3, 333, time.Now(), "pilot3", "glider3", "glider_id3", 20.3, "http://test3.local"})

	// Uncomment this, run the unit-test for the discord webhook to get the message.
	// CheckWebhooks()

	// The discord should only resieve 1 message, containg 3 track IDs.
	// This have to be checked manually.
	// Example output in Discord: Latest timestamp: 333, 3 new tracks are id1,id2,id3.(processing:837.412357ms)

	// Deletes all tracks and webhooks from the database.
	database.DeleteAll()
	databaseTracks.DeleteAll()
}

// Function to test: NewWebhook()
// Test if the insertion of a wrong formatet ID returns correct error.
func Test_NewWebhook_MalformedPost(t *testing.T) {
	CollectionTrack = "TestTracks"
	CollectionWebhook = "TestWebhooks"

	// Connects to a daabase.
	database := mongodb.DatabaseInit(CollectionWebhook)

	// Data to send, is in wrong format.
	postString := "{\"wrong\":\"wrong\"}"

	// Creates a request that is passed to the handler.
	request, _ := http.NewRequest("POST", "/paragliding/api/webhook/new_track/", strings.NewReader(postString))

	// Creates the recorder and router.
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()

	// Tests the function.
	router.HandleFunc("/paragliding/api/webhook/new_track/", NewWebhook).Methods("POST")
	router.ServeHTTP(recorder, request)

	// Check the status code is what we expect (400).
	status := recorder.Code
	if status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	// Check if the content-type is what we expect (application/json).
	content := recorder.HeaderMap.Get("content-type")
	if content != "application/json" {
		t.Errorf("Handler returned wrong content-type: got %s want %s",
			content, "application/json")
	}

	// Check if the handler returns correct data.
	actual := recorder.Body.String()
	expected := "malformed POST request, should be '{\"webhookURL\": \"<url>\"}, [optional: (\"minTriggerValue\": <number>')]"
	if actual != expected {
		t.Errorf("Handler returned wrong error: got %s want %s",
			actual, expected)
	}

	// Deletes all tracks and webhooks from the database.
	database.DeleteAll()
}

// Function to test: NewWebhook()
// Test if the correct error is displayed when duplicate webhook is posted.
func Test_NewWebhook_DuplicateWebhook(t *testing.T) {
	CollectionTrack = "TestTracks"
	CollectionWebhook = "TestWebhooks"

	// Connects to a daabase.
	database := mongodb.DatabaseInit(CollectionWebhook)

	// Inserts a webhook.
	database.InsertWebhook(mongodb.Webhook{"http://test1.local", 3, 0})

	// Try to post a webhook, that allready exists.
	postString := "{ \"webhookURL\": \"http://test1.local\", \"minTriggerValue\": 3 }"

	// Creates a request that is passed to the handler.
	request, _ := http.NewRequest("POST", "/paragliding/api/webhook/new_track/", strings.NewReader(postString))

	// Creates the recorder and router.
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()

	// Tests the function.
	router.HandleFunc("/paragliding/api/webhook/new_track/", NewWebhook).Methods("POST")
	router.ServeHTTP(recorder, request)

	// Check the status code is what we expect (Conflict 409).
	status := recorder.Code
	if status != http.StatusConflict {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusConflict)
	}

	// Check if the handler returns correct data.
	actual := recorder.Body.String()
	expected := "the webhook allready exists"
	if actual != expected {
		t.Errorf("Handler returned wrong error: got %s want %s",
			actual, expected)
	}

	// Deletes all tracks and webhooks from the database.
	database.DeleteAll()
}

// Function to test: NewWebhook()
// Test if the 'minTriggerValue' is set to 1, when not provided.
func Test_NewWebhook_MinTriggerField(t *testing.T) {
	CollectionTrack = "TestTracks"
	CollectionWebhook = "TestWebhooks"

	// Connects to a daabase.
	database := mongodb.DatabaseInit(CollectionWebhook)

	// Data to send.
	postString := "{ \"webhookURL\": \"http://test2.local\" }"

	// Creates a request that is passed to the handler.
	request, _ := http.NewRequest("POST", "/paragliding/api/webhook/new_track/", strings.NewReader(postString))

	// Creates the recorder and router.
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()

	// Tests the function.
	router.HandleFunc("/paragliding/api/webhook/new_track/", NewWebhook).Methods("POST")
	router.ServeHTTP(recorder, request)

	// Check the status code is what we expect (201).
	status := recorder.Code
	if status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Check if the content-type is what we expect (application/json).
	content := recorder.HeaderMap.Get("content-type")
	if content != "application/json" {
		t.Errorf("Handler returned wrong content-type: got %s want %s",
			content, "application/json")
	}

	// The actual returned data.
	id := recorder.Body.String()

	// Get the track from the database.
	expectedTrack, _ := database.FindWebhook(id)

	// Check if the webhook returned has the same url as the one posted.
	if expectedTrack.WebhookURL != "http://test2.local" {
		t.Error("Handler returned wrong ID for webhook")
	}

	// Check if the MinTriggerValue field got updatet.
	if expectedTrack.MinTriggerValue != 1 {
		t.Errorf("Handler did not update the field 'minTriggerValue': got %d want %d",
			expectedTrack.MinTriggerValue, 1)
	}

	// Deletes all webhooks from the database.
	database.DeleteAll()
}

// Function to test: NewWebhook()
// Test if the insertion of a new webhook works, and returns the correct ID.
func Test_NewWebhook(t *testing.T) {
	CollectionTrack = "TestTracks"
	CollectionWebhook = "TestWebhooks"

	// Connects to a daabase.
	database := mongodb.DatabaseInit(CollectionWebhook)

	// Data to send.
	postString := "{ \"webhookURL\": \"http://test1.local\", \"minTriggerValue\": 3 }"

	// Creates a request that is passed to the handler.
	request, _ := http.NewRequest("POST", "/paragliding/api/webhook/new_track/", strings.NewReader(postString))

	// Creates the recorder and router.
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()

	// Tests the function.
	router.HandleFunc("/paragliding/api/webhook/new_track/", NewWebhook).Methods("POST")
	router.ServeHTTP(recorder, request)

	// Check the status code is what we expect (201).
	status := recorder.Code
	if status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Check if the content-type is what we expect (application/json).
	content := recorder.HeaderMap.Get("content-type")
	if content != "application/json" {
		t.Errorf("Handler returned wrong content-type: got %s want %s",
			content, "application/json")
	}

	// The actual returned data.
	id := recorder.Body.String()

	// Get the track from the database.
	expectedTrack, _ := database.FindWebhook(id)

	// Check if the webhook returned has the same url as the one posted.
	if expectedTrack.WebhookURL != "http://test1.local" {
		t.Error("Handler returned wrong ID for webhook")
	}

	// Deletes all webhooks from the database.
	database.DeleteAll()
}

// Function to test: HandleWebhooks()
// Test if correct error code is returned when wrong method is used.
func Test_HandleWebhooks(t *testing.T) {
	// Creates a request that is passed to the handler.
	request, _ := http.NewRequest("POST", "/paragliding/api/webhook/new_track/111", strings.NewReader("hello"))

	// Creates the recorder and router.
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()

	// Tests the function.
	router.HandleFunc("/paragliding/api/webhook/new_track/111", HandleWebhooks).Methods("POST")
	router.ServeHTTP(recorder, request)

	// Check the status code is what we expect (400).
	status := recorder.Code
	if status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

// Function to test: getWebhookInfo().
// Test if the information returned when the webhook does not exist.
func Test_getWebhookInfo_NoWebhook(t *testing.T) {
	// Creates the request.
	request, _ := http.NewRequest("GET", "/paragliding/api/webhook/new_track/11", nil)

	// Creates the recorder and router.
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()

	// Tests the function.
	router.HandleFunc("/paragliding/api/webhook/new_track/11", getWebhookInfo).Methods("GET")
	router.ServeHTTP(recorder, request)

	// Check the status code is what we expect (404).
	status := recorder.Code
	if status != http.StatusNotFound {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

// Function to test: getWebhookInfo().
// Test if the information returned is the same as in the DB.
func Test_getWebhookInfo(t *testing.T) {
	CollectionWebhook = "TestWebhooks"

	// Connects to a daabase.
	database := mongodb.DatabaseInit(CollectionWebhook)

	// Creates the new webhook.
	newHook := mongodb.Webhook{"http://test1.local", 3, 0}

	// Inserts the webhook, and get its ID.
	id, _ := database.InsertWebhook(newHook)

	// Creates the request.
	requestPath := string("/paragliding/api/webhook/new_track/" + id)
	request, _ := http.NewRequest("GET", requestPath, nil)

	// Creates the recorder and router.
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()

	// Tests the function.
	router.HandleFunc(requestPath, HandleWebhooks).Methods("GET")
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

	// Expected return data.
	toJSON, _ := json.Marshal(newHook)
	expected := string(toJSON)

	// The actual returned data.
	actual := recorder.Body.String()

	// Check if the expected is the same as the actual returned data.
	if actual != expected {
		t.Errorf("Handler returned wrong data: got %s want %s",
			actual, expected)
	}

	// Deletes all webhooks from the database.
	database.DeleteAll()
}

// Function to test: deleteWebhook().
// Test if the information returned when the webhook does not exist.
func Test_deleteWebhook_NoWebhook(t *testing.T) {
	// Creates the request.
	request, _ := http.NewRequest("DELETE", "/paragliding/api/webhook/new_track/11", nil)

	// Creates the recorder and router.
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()

	// Tests the function.
	router.HandleFunc("/paragliding/api/webhook/new_track/11", HandleWebhooks).Methods("DELETE")
	router.ServeHTTP(recorder, request)

	// Check the status code is what we expect (404).
	status := recorder.Code
	if status != http.StatusNotFound {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

// Function to test: deleteWebhook()
// Test if correct data is returned and if the webhook got deleted from the DB.
func Test_deleteWebhook(t *testing.T) {
	CollectionWebhook = "TestWebhooks"

	// Connects to a daabase.
	database := mongodb.DatabaseInit(CollectionWebhook)

	// Creates the new webhook.
	newHook := mongodb.Webhook{"http://test1.local", 3, 0}

	// Inserts the webhook, and get its ID.
	id, _ := database.InsertWebhook(newHook)

	// Creates the request.
	requestPath := string("/paragliding/api/webhook/new_track/" + id)
	request, _ := http.NewRequest("DELETE", requestPath, nil)

	// Creates the recorder and router.
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()

	// Tests the function.
	router.HandleFunc(requestPath, HandleWebhooks).Methods("DELETE")
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

	// Expected return data.
	toJSON, _ := json.Marshal(newHook)
	expected := string(toJSON)

	// The actual returned data.
	actual := recorder.Body.String()

	// Check if the expected is the same as the actual returned data.
	if actual != expected {
		t.Errorf("Handler returned wrong data: got %s want %s",
			actual, expected)
	}

	// Check if the webhook got deleted from the DB.
	expectedError := "not found"
	_, err := database.DeleteWebhook(id)

	// Check if any errors got returned.
	if err != nil {
		// Check if the correct error got returned.
		if err.Error() != expectedError {
			t.Errorf("Handler returned wrong error: got %s want %s",
				err.Error(), expectedError)
		}
	}
}
