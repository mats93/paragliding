/*
  File: webhookDatabase_test.go
  Contains unit tests for webhookDatabase.go
*/

package mongodb

import (
	"testing"
	"time"
)

// Method to test: InsertWebhook().
// Test if insertion of a webhook works, and correct errors are returned when not.
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

// Method to test: InvokeWebhooks().
// Test the method that invokes the webhooks for correct updates in the DB and returns.
func Test_InvokeWebhooks(t *testing.T) {
	// Connects to database and insert webhooks to test.
	database := DatabaseInit("TestWebhooks")
	id, _ := database.InsertWebhook(Webhook{"http://webhook.local", 2, 0})
	database.InsertWebhook(Webhook{"http://webhook2.local", 99, 0})
	defer MDB.Session.Close()

	// Connects to database and inserts a track.
	databaseTracks := DatabaseInit("TestTracks")
	databaseTracks.Insert(Track{1, 111, time.Now(), "pilot1", "glider1", "glider_id1", 20.1, "http://test1.local"})
	defer MDB.Session.Close()

	// Test the method.
	database = DatabaseInit("TestWebhooks")
	hooksEmpty, _ := database.InvokeWebhooks()

	// One track has been insertet. hooks should not contain any Webhooks.
	if hooksEmpty != nil {
		t.Error("Method did return a webhook, when no webhooks should be returned")
	}

	// Check if the 'NumberOfNewInserts' has neen updatet.
	dbHook, _ := database.FindWebhook(id)
	if dbHook.NumberOfNewInserts != 1 {
		t.Error("Method did not update the field 'numberOfNewInserts'")
	}

	// Check if the correct webhook is retrieved.
	if dbHook.WebhookURL != "http://webhook.local" {
		t.Errorf("Method returned wrong Webhook: got %s want %s",
			dbHook.WebhookURL, "http://webhook.local")
	}

	// Connects to database and inserts a new track.
	databaseTracks = DatabaseInit("TestTracks")
	databaseTracks.Insert(Track{2, 222, time.Now(), "pilot2", "glider2", "glider_id2", 20.2, "http://test2.local"})
	defer MDB.Session.Close()

	// Test the method again.
	database = DatabaseInit("TestWebhooks")
	hooks, _ := database.InvokeWebhooks()

	// Two track has been insertet, hooks should therefore contain the Webhook.
	if hooks == nil {
		t.Errorf("Method did not return a webhook when expected to")
	}

	// Deletes all webhooks from the database.
	database.DeleteAll()
	databaseTracks.DeleteAll()

	// Closes the database session.
	defer MDB.Session.Close()
}

// Method to test: FindWebhook().
// Test if the correct error messages and webhook is retrived.
func Test_FindWebhook(t *testing.T) {
	// Connects to database and insert webhooks to test.
	database := DatabaseInit("TestWebhooks")
	database.InsertWebhook(Webhook{"http://webhook.local", 1, 0})
	id, _ := database.InsertWebhook(Webhook{"http://webhook2.local", 2, 0})
	database.InsertWebhook(Webhook{"http://webhook3.local", 3, 0})

	// Expected error message.
	errorMessage := "not found"

	// Check if the correct error is returned when the provided ID does not exist.
	_, err := database.FindWebhook("1")
	if err == nil {
		t.Errorf("Method did not return error when expected to")

		if err.Error() != errorMessage {
			t.Errorf("Method returned wrong error: got \"%s\" want \"%s\"",
				err.Error(), errorMessage)
		}
	}

	// Check if the correct webhook is returned when provided with correct ID.
	hook, _ := database.FindWebhook(id)
	expectedURL := "http://webhook2.local"
	if hook.WebhookURL != expectedURL {
		t.Errorf("Method returned wrong hook (URL): got %s want %s",
			hook.WebhookURL, expectedURL)
	}

	// Deletes all webhooks from the database.
	database.DeleteAll()

	// Closes the database session.
	defer MDB.Session.Close()
}

// Method to test: DeleteWebhook().
// Test if the correct error messages and webhook is retrived.
func Test_DeleteWebhook(t *testing.T) {
	// Connects to database and insert webhooks to test.
	database := DatabaseInit("TestWebhooks")
	database.InsertWebhook(Webhook{"http://webhook.local", 1, 0})
	id, _ := database.InsertWebhook(Webhook{"http://webhook2.local", 2, 0})
	database.InsertWebhook(Webhook{"http://webhook3.local", 3, 0})

	// Expected error message.
	errorMessage := "not found"

	// Check if the error message is correct when deleting
	// a non-existing webhook.
	_, err := database.DeleteWebhook("1")

	if err == nil {
		t.Errorf("Method did not return error when expected to")

		if err.Error() != errorMessage {
			t.Errorf("Method returned wrong error: got \"%s\" want \"%s\"",
				err.Error(), errorMessage)
		}
	}

	// Check the if the returned webhook is correct
	// when deleting a webhook with correct ID.
	hook, _ := database.DeleteWebhook(id)
	expectedURL := "http://webhook2.local"
	if hook.WebhookURL != expectedURL {
		t.Errorf("Method returned wrong hook (URL): got %s want %s",
			hook.WebhookURL, expectedURL)
	}

	// Check if the deleted webhook is removed from the DB.
	_, err = database.FindWebhook(id)
	if err == nil {
		t.Error("Method did not delete the webhook")
	}

	// Deletes all webhooks from the database.
	database.DeleteAll()

	// Closes the database session.
	defer MDB.Session.Close()
}
