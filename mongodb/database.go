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

// COLLECTION is a MongoDB collection.
const COLLECTION = "Tracks"

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
type MongoDB struct {
	Server   string
	Database string
	Username string
	Password string
}

// Name of the Database that is connected.
var MDB *mgo.Database

// Connect to the database.
func (m *MongoDB) Connect() {
	// Creates the Database URL.
	dbURL := string(m.Username + ":" + m.Password + "@" + m.Server + "/" + m.Database)

	// Starts the session.
	// An error in the mgo.Dial wil create a Panic(err) in the mgo package.
	session, _ := mgo.Dial(dbURL)

	// Saves the name of the Database that was connected.
	MDB = session.DB(m.Database)
}

// Find all entries in the collection.
func (m *MongoDB) FindAll() ([]Track, error) {
	var results []Track

	// Find all tracks in the collection.
	err := MDB.C(COLLECTION).Find(bson.M{}).All(&results)

	// Returns the struct, and error if any.
	return results, err
}

// Find entry by ID.
func (m *MongoDB) FindByID(id int) ([]Track, error) {
	var result []Track

	// Find track with given 'id'.
	err := MDB.C(COLLECTION).Find(bson.M{"id": id}).All(&result)

	// Generate error if track with given ID was not found.
	if result == nil {
		err = errors.New("not found")
	}
	// Returns the struct, and error if any.
	return result, err
}

// Insert a new Struct into the database.
func (m *MongoDB) Insert(t Track) error {
	err := MDB.C(COLLECTION).Insert(&t)
	return err
}

// Get a count of all tracks in the database.
func (m *MongoDB) GetCount() (int, error) {
	count, err := MDB.C(COLLECTION).Count()
	if err != nil {
		return 0, err
	}
	// Returns the count of all tracks and no error.
	return count, nil
}

// Deletes all entries in the database collection.
func (m *MongoDB) DeleteAllTracks() error {
	_, err := MDB.C(COLLECTION).RemoveAll(bson.M{})
	return err
}

// DatabaseInit Initialises the database, and connects to it.
func DatabaseInit() MongoDB {
	database := MongoDB{
		"ds233763.mlab.com:33763",
		"paragliding_db",
		"paragliderAPI",
		"6oLKQOFcxMDCZyd",
	}
	// Connets to the database and returns the struct.
	database.Connect()
	return database
}

// TestDb is just for testing.
func TestDb() {

	// Database settings.
	database := MongoDB{
		"ds233763.mlab.com:33763",
		"paragliding_db",
		"paragliderAPI",
		"6oLKQOFcxMDCZyd",
	}

	// Connects to the database.
	database.Connect()

	/* Inserts a track to the database.
	database.Insert(Track{1, time.Now(), "pilot", "glider", "glider_id", 20.4, "http://test.test"})
	database.Insert(Track{2, time.Now(), "pilot", "glider", "glider_id", 20.4, "http://test.test"})
	database.Insert(Track{3, time.Now(), "pilot", "glider", "glider_id", 20.4, "http://test.test"})
	database.Insert(Track{4, time.Now(), "pilot", "glider", "glider_id", 20.4, "http://test.test"})
	*/

	// Deletes all tracks from the database.
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

	// Gets track with ID from the database.
	trackID, err := database.FindByID(1)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("\nBy id:", trackID)
	}

	// Gets the count of all tracks in the database.
	count, err := database.GetCount()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Det er %d tracks i databasen \n", count)
	}
	// Closes the database session.
	defer MDB.Session.Close()
}
