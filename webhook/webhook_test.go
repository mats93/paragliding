/*
	File: webhook_test.go
  Contains unit tests for webhook.go
*/

package webhook

import (
	"testing"
	"time"

	"github.com/mats93/paragliding/mongodb"
)

// Test if the webhooks get the payload.
func Test_CheckWebhooks(t *testing.T) {
	CollectionTrack = "Tests"
	CollectionWebhook = "TestWebhooks"

	// Adds some webhooks.
	database := mongodb.DatabaseInit(CollectionWebhook)
	database.InsertWebhook(mongodb.Webhook{
		"https://discordapp.com/api/webhooks/504733605344313354/8sLrUSTCJQxB-8BAcRBH27T3jm8xWtCv1DjpjdbJnhqWqNsjxbatT_EWFsQtPnKzuuCf", 3, 0})
	database.InsertWebhook(mongodb.Webhook{"http://test2.local", 2, 0})
	database.InsertWebhook(mongodb.Webhook{"http://test3.local", 1, 0})

	databaseTracks := mongodb.DatabaseInit(CollectionTrack)

	// Add 3 tracks to the DB.
	databaseTracks.Insert(mongodb.Track{1, 111, time.Now(), "pilot1", "glider1", "glider_id1", 20.1, "http://test1.local"})
	CheckWebhooks()
	databaseTracks.Insert(mongodb.Track{2, 222, time.Now(), "pilot2", "glider2", "glider_id2", 20.2, "http://test2.local"})
	CheckWebhooks()
	databaseTracks.Insert(mongodb.Track{3, 333, time.Now(), "pilot3", "glider3", "glider_id3", 20.3, "http://test3.local"})
	CheckWebhooks()

	// Deletes all tracks and webhooks from the database.
	database.DeleteAll()
	databaseTracks.DeleteAll()
}
