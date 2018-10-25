/*
	File: webhook.go
  Contains functions used by API calls to the "Webhook paths".
*/

package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mats93/paragliding/mongodb"
)

// Message to be sent to subscrber.
type notifyMessage struct {
	URL           string        `json:"-"`
	HumanReadable string        `json:"-"`
	TimeLatest    int64         `json:"t_latest"`
	Tracks        []int         `json:"tracks"`
	Processing    time.Duration `json:"processing"`
}

// The DB collections to use, gets injected from main or test.
var CollectionTrack string
var CollectionWebhook string

// Checks if the registrated webhooks need to notify the subscrber.
// This function should be called everytime a track is added.
func CheckWebhooks() {
	// Start time of the request.
	start := time.Now()

	// Connects to the database, uses the Webhook collection.
	dbWebhooks := mongodb.DatabaseInit(CollectionWebhook)

	// Adds 1 to the 'newInserts' field, and returns a slice
	// of all webhooks that need to be notified.
	webhooks, err := dbWebhooks.InvokeWebhooks()
	if err != nil {
		// Logs the error.
		log.Fatal(err)
	}
	if webhooks != nil {
		// Creates a new db session against track collection.
		dbTracks := mongodb.DatabaseInit(CollectionTrack)

		// Getts all tracks from the database.
		tracks, err := dbTracks.FindAll()
		if err != nil {
			// Logs the error.
			log.Fatal(err)
		} else {
			var messages []notifyMessage
			var ids []int

			// The last index in the tracks slice.
			lenTracks := len(tracks) - 1

			// Loops through all webhooks.
			for i := 0; i < len(webhooks); i++ {

				// Loop through all tracks that should be in the message, and adds them.
				for j := 0; j < webhooks[i].NumberOfNewInserts+1; j++ {
					ids = append(ids, tracks[j].ID)
				}

				// Creates the new message and adds it to a slice.
				messages = append(messages, notifyMessage{
					webhooks[i].WebhookURL, "", tracks[lenTracks].Timestamp, ids, time.Since(start)})

				// Removes content from the slice.
				ids = nil
			}

			// Formats the messages to a human-readable format.
			formatIDs := ""
			formatMessage := ""

			// Loops through all messages.
			for i := 0; i < len(messages); i++ {
				formatIDs = ""
				// Loops through all IDs.
				for j := 0; j < len(messages[i].Tracks); j++ {
					// Dont add "," on last id.
					if j == len(messages[i].Tracks)-1 {
						formatIDs += fmt.Sprintf("id%v", messages[i].Tracks[j])
					} else {
						formatIDs += fmt.Sprintf("id%v,", messages[i].Tracks[j])
					}
				}
				// Formats the message.
				formatMessage = fmt.Sprintf("Latest timestamp: %d, %d new tracks are %s.(processing:%v)",
					messages[i].TimeLatest, len(messages[i].Tracks), formatIDs, messages[i].Processing)

				// Converts formatet message to json, uses 'content' to support Discord webhooks.
				jsonMessage := map[string]string{"content": formatMessage}
				body, err := json.Marshal(jsonMessage)
				if err != nil {
					log.Fatal(err)
				}

				// Send the message.
				resp, err := http.Post(messages[i].URL, "application/json", bytes.NewBuffer(body))
				if err == nil {
					defer resp.Body.Close()
				}
			}
		}
	}
}

// NewWebhook - POST: Registrates a new webhook.
// Output: application/json
func NewWebhook(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method.
	if r.Method != "POST" {
		// If the request is not POST.
		// Sets header content-type to application/json and status code to 400 (Bad request).
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
	} else {
		// The request is POST.
		// Connects the the database.
		database := mongodb.DatabaseInit(CollectionWebhook)
		var newWebhook mongodb.Webhook

		// Decodes the json url and converts it to a struct.
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&newWebhook)

		// Check for error or if the 'URL' section is empty.
		if err != nil || newWebhook.WebhookURL == "" {
			// The decoding failed.
			// Sets header status code to 400 "Bad request", and returns error message.
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("malformed POST request, should be '{\"webhookURL\": \"<url>\"}, [optional: (\"minTriggerValue\": <number>')]"))
		} else {
			// The decoding was sucessful.

			// Check if the optional field 'minTriggerValue' is set.
			if newWebhook.MinTriggerValue == 0 {
				// If not, set it to 1 (default).
				newWebhook.MinTriggerValue = 1
			}

			// Adds the new webhook to the db.
			id, err := database.InsertWebhook(newWebhook)
			if err != nil {
				if err.Error() == "the webhook allready exists" {
					// The request is valid but it exists.
					// but the webhook is allready registrated, output error code 409 and error.
					w.WriteHeader(http.StatusConflict)
					w.Write([]byte(err.Error()))
				} else {
					// Sets header status code to 500 "Internal server error".
					w.WriteHeader(http.StatusInternalServerError)
					log.Fatal(err)
				}
			} else {
				// The webhook was created, returns the ID.
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(id))
			}
		}
	}
}

// Gets information about a registrated webhook.
func getWebhookInfo(w http.ResponseWriter, r *http.Request) {
	// Connects the the database.
	database := mongodb.DatabaseInit(CollectionWebhook)
	var id string

	// Gets the ID from the URL and converts it to a string.
	fmt.Sscanf(r.URL.Path, "/paragliding/api/webhook/new_track/%s", &id)

	// Gets the webhook from the DB, with the given ID.
	hook, err := database.FindWebhook(id)
	if err != nil {
		// Error: A webhook with the given ID does not excist.
		// Sets the header code to 404 (Not found).
		w.WriteHeader(http.StatusNotFound)

	} else {
		// The webhook was retrieved from the DB.
		// Converts the struct to json and outputs it.
		json, err := json.Marshal(hook)
		if err != nil {
			// Error: Sets header status code to 500 "Internal server error" and logs the error.
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
		} else {
			// OK: Sets header content-type to application/json and status code to 200 (OK).
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			// Returns the webhook as json.
			w.Write([]byte(json))
		}
	}
}

// Deletes a webhook.
func deleteWebhook(w http.ResponseWriter, r *http.Request) {
	// Connects the the database.
	database := mongodb.DatabaseInit(CollectionWebhook)
	var id string

	// Gets the ID from the URL and converts it to a string.
	fmt.Sscanf(r.URL.Path, "/paragliding/api/webhook/new_track/%s", &id)

	// Delete the webhook from the DB, with the given ID.
	hook, err := database.DeleteWebhook(id)
	if err != nil {
		// Error: A webhook with the given ID does not excist.
		// Sets the header code to 404 (Not found).
		w.WriteHeader(http.StatusNotFound)

	} else {
		// The webhook was retrieved from the DB.
		// Converts the struct to json and outputs it.
		json, err := json.Marshal(hook)
		if err != nil {
			// Error: Sets header status code to 500 "Internal server error" and logs the error.
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
		} else {
			// OK: Sets header content-type to application/json and status code to 200 (OK).
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			// Returns the webhook as json.
			w.Write([]byte(json))
		}
	}
}

// WebhookHandler - GET:    Get information about a given webhook.
// WebhookHandler - DELETE: Delete a given webhook.
// Output: application/json
func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	// Calls functions to handle the GET and POST requests.
	switch r.Method {
	case "GET":
		getWebhookInfo(w, r)

	case "DELETE":
		deleteWebhook(w, r)

	default:
		// Wrong method.
		// Sets status code to 400 (Bad request).
		w.WriteHeader(http.StatusBadRequest)
	}
}
