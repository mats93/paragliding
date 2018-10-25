/*
	File: webhookDatabase.go
  Handles the mongoDB operations for webhooks
*/

package mongodb

import (
	"encoding/hex"
	"errors"

	"github.com/globalsign/mgo/bson"
)

// Webhook struct.
type Webhook struct {
	WebhookURL         string `bson:"webhookURL"         json:"webhookURL"`
	MinTriggerValue    int    `bson:"minTriggerValue"    json:"minTriggerValue"`
	NumberOfNewInserts int    `bson:"numberOfNewInserts" json:"-"`
}

// ID of MongoDB webhook object.
type bsonID struct {
	ID bson.ObjectId `json:"id" bson:"_id,omitempty"`
}

// InsertWebhook inserts a new webhook to the database.
func (m *MongoDB) InsertWebhook(hook Webhook) (string, error) {
	var id bsonID

	// Check if the webhook allready exists.
	count, err := MDB.C(m.Collection).Find(bson.M{"webhookURL": hook.WebhookURL}).Count()

	if err != nil {
		// There was an error in searching the DB.
		return "", err
	}
	// The DB search was sucessful.

	if count == 0 {
		// The webhook does not exist.

		// Inserts the webhook to the database.
		err = MDB.C(m.Collection).Insert(&hook)

		if err != nil {
			// There was an error in inserting the webhook to the DB.
			return "", err
		}
		// The insertion was sucessful.

		// Queries the DB for the URL.
		err := MDB.C(m.Collection).Find(bson.M{"webhookURL": hook.WebhookURL}).One(&id)
		if err != nil {
			// There was an error in finding the ID.
			return "", err
		}
		// Returns the bson ID and no error.
		return id.ID.Hex(), nil

	}
	// The webhook allready exists, return error.
	return "", errors.New("the webhook allready exists")
}

// InvokeWebhooks invokes all webhooks that meet the criteria.
// This method is called everytime a new track is inserted.
func (m *MongoDB) InvokeWebhooks() ([]Webhook, error) {
	var results []Webhook
	var returnedHooks []Webhook

	// Find all webhook in the collection.
	err := MDB.C(m.Collection).Find(bson.M{}).All(&results)

	if err != nil {
		// Logs error.
		return nil, err
	}
	// Loops through all webhooks.
	for i := 0; i < len(results); i++ {

		// Updates all 'NumberOfNewInserts' by 1.
		err = MDB.C(m.Collection).Update(bson.M{"webhookURL": results[i].WebhookURL},
			bson.M{"$set": bson.M{"numberOfNewInserts": results[i].NumberOfNewInserts + 1}})
		if err != nil {
			// Returns error.
			return nil, err
		}

		// Check if the subscriber should be notified.
		if results[i].NumberOfNewInserts+1 >= results[i].MinTriggerValue {
			// Adds the webhook that should be notified to a new slice.
			returnedHooks = append(returnedHooks, results[i])

			// Sets the 'newInserts' value back to 0.
			err = MDB.C(m.Collection).Update(bson.M{"webhookURL": results[i].WebhookURL},
				bson.M{"$set": bson.M{"numberOfNewInserts": 0}})
			if err != nil {
				// Returns error.
				return nil, err
			}
		}
	}
	// Returns the slice of Webhooks that sohuld be notified.
	return returnedHooks, nil
}

// IsObjectIDHex checks if an ID can be a mongodb ID.
func IsObjectIDHex(s string) bool {
	if len(s) != 24 {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil
}

// FindWebhook finds a webhook by ID.
func (m *MongoDB) FindWebhook(id string) (Webhook, error) {
	var result []Webhook
	errorMessage := "not found"

	// Check if the ID can be a mongodb ID.
	if IsObjectIDHex(id) {
		// True, find track with given 'id'.
		err := MDB.C(m.Collection).Find(bson.M{"_id": bson.ObjectIdHex(id)}).All(&result)

		// Generate error if webhook with given ID was not found.
		if result == nil {
			// Returns empty struct and the error.
			err = errors.New(errorMessage)
			return Webhook{}, err
		}

		// Returns the struct, and error if any.
		return result[0], err
	}
	// False, return empty struct and error.
	err := errors.New(errorMessage)
	return Webhook{}, err
}

// DeleteWebhook deletes a webhook with a given ID.
func (m *MongoDB) DeleteWebhook(id string) (Webhook, error) {
	var result []Webhook
	errorMessage := "not found"

	// Check if the ID can be a mongodb ID.
	if IsObjectIDHex(id) {
		// True, get the webhook that should be deleted.
		err := MDB.C(m.Collection).Find(bson.M{"_id": bson.ObjectIdHex(id)}).All(&result)

		// Generate error if webhook with given ID was not found.
		if result == nil {
			// Returns empty struct and the error.
			err = errors.New(errorMessage)
			return Webhook{}, err
		}

		// Delete the webhook from the database.
		err = MDB.C(m.Collection).Remove(bson.M{"_id": bson.ObjectIdHex(id)})
		if err != nil {
			// Error in removing document, returns the error.
			return Webhook{}, err
		}
		// Returns the deleted webhook and no error.
		return result[0], nil
	}
	// False, return empty strcut and error.
	err := errors.New(errorMessage)
	return Webhook{}, err
}
