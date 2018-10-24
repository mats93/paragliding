/*
	File: webhookDatabase.go
  Handles the mongoDB operations for webhooks
*/

package mongodb

import (
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

// Inserts a new webhook to the database.
func (m *MongoDB) InsertWebhook(hook Webhook) (string, error) {
	var bsonId bsonID

	// Check if the webhook allready exists.
	count, err := MDB.C(m.Collection).Find(bson.M{"webhookURL": hook.WebhookURL}).Count()

	if err != nil {
		// There was an error in searching the DB.
		return "", err
	} else {
		// The DB search was sucessful.

		if count == 0 {
			// The webhook does not exist.

			// Inserts the webhook to the database.
			err = MDB.C(m.Collection).Insert(&hook)

			if err != nil {
				// There was an error in inserting the webhook to the DB.
				return "", err
			} else {
				// The insertion was sucessful.

				// Queries the DB for the URL.
				err := MDB.C(m.Collection).Find(bson.M{"webhookURL": hook.WebhookURL}).One(&bsonId)
				if err != nil {
					// There was an error in finding the ID.
					return "", err
				}
				// Returns the bson ID and no error.
				return bsonId.ID.Hex(), nil
			}
		} else {
			// The webhook allready exists, return error.
			return "", errors.New("the webhook allready exists")
		}
	}
}

// Invokes webhooks that meet the criteria.
// This method is called everytime a new track is inserted.
func (m *MongoDB) InvokeWebhooks() ([]Webhook, error) {
	var results []Webhook
	var returnedHooks []Webhook

	// Find all webhook in the collection.
	err := MDB.C(m.Collection).Find(bson.M{}).All(&results)

	if err != nil {
		// Logs error.
		return nil, err
	} else {
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
}

// Finds a webhook by ID.
func (m *MongoDB) FindWebhook(id int) (Webhook, error) {
	var hook Webhook

	return hook, nil
}

// Deletes a webhook with a given ID.
func (m *MongoDB) DeleteWebhook(id int) (Webhook, error) {
	var hook Webhook

	return hook, nil
}
