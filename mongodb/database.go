/*
	File: database.go
  Handles the mongoDB operations.
*/

package mongodb

import (
	"errors"
	"fmt"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Track is the metadata about the track that will be stored in the database.
type Track struct {
	ID          int       `json:"-"`
	HDate       time.Time `bson:"H_date"        json:"H_date"`
	Pilot       string    `bson:"pilot"         json:"pilot"`
	Glider      string    `bson:"glider"        json:"glider"`
	GliderID    string    `bson:"glider_id"     json:"glider_id"`
	TrackLength float64   `bson:"track_length"  json:"track_length"`
	TrackSrcURL string    `bson:"track_src_url" json:"track_src_url"`
}

// DatabaseMGO holds the database information.
type databaseMGO struct {
	Server   string
	Database string
	Username string
	Password string
}

// Name of the Database.
var db *mgo.Database

// COLLECTION is a mgo database collection.
const COLLECTION = "test"

// Connect to the database.
func (m *databaseMGO) Connect() {
	// Creates the Database URL.
	dbURL := string(m.Username + ":" + m.Password + "@" + m.Server + "/" + m.Database)

	// Starts the session.
	// An error in the mgo.Dial wil create a Panic(err) in the mgo package.
	session, _ := mgo.Dial(dbURL)

	// Saves the name of the Database that was connected.
	db = session.DB(m.Database)
}

// Database Queries:

// Find all entries in the collection.
func (m *databaseMGO) FindAll() ([]Track, error) {
	var results []Track

	// Find all tracks in the collection.
	err := db.C(COLLECTION).Find(bson.M{}).All(&results)

	// Returns the struct, and error if any.
	return results, err
}

// Find entry by ID.
func (m *databaseMGO) FindByID(id int) ([]Track, error) {
	var result []Track

	// Find track with given 'id'.
	err := db.C(COLLECTION).Find(bson.M{"id": id}).All(&result)

	// Generate error if track with given ID was not found.
	if result == nil {
		err = errors.New("not found")
	}

	// Returns the struct, and error if any.
	return result, err
}

// Insert a new Struct into the database.
func (m *databaseMGO) Insert(t Track) error {
	err := db.C(COLLECTION).Insert(&t)
	return err
}

// Get a count of all tracks in the database.

// Delete all entries in the database collection.
func (m *databaseMGO) DeleteAllTracks() {
	_, err := db.C(COLLECTION).RemoveAll(bson.M{})
	if err != nil {
		fmt.Println(err)
	}
}

// TestDb is just for testing.
func TestDb() {

	// Database settings.
	database := databaseMGO{"ds233763.mlab.com:33763", "paragliding_db", "admin", "passord1"}

	// Connects to the database.
	database.Connect()

	/* Inserts a track to the database.
	database.Insert(Track{1, time.Now(), "pilot", "glider", "glider_id", 20.4, "http://test.test"})
	database.Insert(Track{2, time.Now(), "pilot", "glider", "glider_id", 20.4, "http://test.test"})
	database.Insert(Track{3, time.Now(), "pilot", "glider", "glider_id", 20.4, "http://test.test"})
	database.Insert(Track{4, time.Now(), "pilot", "glider", "glider_id", 20.4, "http://test.test"})
	*/

	// database.DeleteAllTracks()

	// Gets all tracks from the database.
	results, err := database.FindAll()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Results:", results)
		for i := 0; i < len(results); i++ {
			fmt.Println(results[i].ID)
		}
	}

	// Gets track with ID 5 from the database.
	test, err := database.FindByID(1)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("\nBy id:", test)
	}

	// Closes the database session.
	defer db.Session.Close()
}
