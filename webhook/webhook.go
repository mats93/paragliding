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

// Collection is the MongoDB collection to use. Gets injected from main or test.
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

}
