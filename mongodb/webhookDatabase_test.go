/*
  File: webhookDatabase_test.go
  Contains unit tests for webhookDatabase.go
*/

package mongodb

import (
	"fmt"
	"testing"
	"time"
)

func Test_InsertWebhook(t *testing.T) {
	// Connects to the database.
	database := DatabaseInit("TestWebhooks")

	// Creates a webhook to test.
	hook := Webhook{"http://test.com", 1, 0}

	// Check if insertions of new unused Webhook works.
	id, err := database.InsertWebhook(hook)

	// Check for unexpected errors.
	if err != nil {
		t.Errorf("Method generated unexpected error: %v", err)
	}

	// Check if an ID was returned.
	if id == "" {
		t.Errorf("Method returned empty ID, when no webhooks allready is in the DB")
	}

	// Check for correct error when a duplicate webhook is inserted.
	expectedError := "the webhook allready exists"
	_, err = database.InsertWebhook(hook)
	if err == nil {
		// Got no error.
		t.Error("Method did not return error when expected to")
	} else {
		// Got error, check if the error is correct.
		actualError := err.Error()

		if actualError != expectedError {
			t.Errorf("Method returned wrong error when duplicate webhook is inserted: got %s want %s",
				actualError, expectedError)
		}
	}

	// Deletes all webhooks from the database.
	database.DeleteAll()

	// Closes the database session.
	defer MDB.Session.Close()
}

func Test_InvokeWebhooks(t *testing.T) {
	databaseTracks := DatabaseInit("Tests")
	databaseTracks.Insert(Track{1, 111, time.Now(), "pilot1", "glider1", "glider_id1", 20.1, "http://test1.test"})
	databaseTracks.Insert(Track{2, 222, time.Now(), "pilot2", "glider2", "glider_id2", 20.2, "http://test2.test"})
	databaseTracks.Insert(Track{3, 333, time.Now(), "pilot3", "glider3", "glider_id3", 20.3, "http://test3.test"})
	defer MDB.Session.Close()

	// Connects to the database.
	database := DatabaseInit("TestWebhooks")

	// Check if insertions of new unused Webhook works.
	fmt.Println(database.InsertWebhook(Webhook{"http://test1.local", 1, 0}))
	fmt.Println(database.InsertWebhook(Webhook{"http://test2.local", 2, 1}))
	fmt.Println(database.InsertWebhook(Webhook{"http://test3.local", 3, 2}))

	//
	hooks, err := database.InvokeWebhooks()
	fmt.Println(hooks)
	fmt.Println(err)

	// Deletes all webhooks from the database.
	database.DeleteAll()
	databaseTracks.DeleteAll()

	// Closes the database session.
	defer MDB.Session.Close()
}
