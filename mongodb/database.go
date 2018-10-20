/*
	File: database.go
  Handles the mongoDB operations.
*/

package mongodb

import (
	"errors"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Track is the metadata about the track that will be stored in the database.
type Track struct {
	ID          int       `json:"-"`
	Timestamp   int64     `json:"-"`
	HDate       time.Time `bson:"H_date"        json:"H_date"`
	Pilot       string    `bson:"pilot"         json:"pilot"`
	Glider      string    `bson:"glider"        json:"glider"`
	GliderID    string    `bson:"glider_id"     json:"glider_id"`
	TrackLength float64   `bson:"track_length"  json:"track_length"`
	TrackSrcURL string    `bson:"track_src_url" json:"track_src_url"`
}

// DatabaseMGO holds the database information.
type MongoDB struct {
	Server     string
	Database   string
	Collection string
	Username   string
	Password   string
}

// The database that connected.
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

// Insert a new Struct into the database.
func (m *MongoDB) Insert(t Track) error {
	err := MDB.C(m.Collection).Insert(&t)
	return err
}

// Deletes all entries in the database collection.
func (m *MongoDB) DeleteAllTracks() error {
	_, err := MDB.C(m.Collection).RemoveAll(bson.M{})
	return err
}

// Find all entries in the collection.
func (m *MongoDB) FindAll() ([]Track, error) {
	var results []Track

	// Find all tracks in the collection.
	err := MDB.C(m.Collection).Find(bson.M{}).All(&results)

	// Returns the struct, and error if any.
	return results, err
}

// Find entry by ID.
func (m *MongoDB) FindByID(id int) ([]Track, error) {
	var result []Track

	// Find track with given 'id'.
	err := MDB.C(m.Collection).Find(bson.M{"id": id}).All(&result)

	// Generate error if track with given ID was not found.
	if result == nil {
		err = errors.New("not found")
	}
	// Returns the struct, and error if any.
	return result, err
}

// Get a count of all tracks in the database.
func (m *MongoDB) GetCount() (int, error) {
	count, err := MDB.C(m.Collection).Count()
	if err != nil {
		return 0, err
	}
	// Returns the count of all tracks and no error.
	return count, nil
}

// Returns a new ID that wil be used in the Track.
func (m *MongoDB) GetNewID() int {
	// For readability, mongoDB`s ID wil not be used.
	// The new ID is simply the Count of all objects in the db + 1.
	count, _ := MDB.C(m.Collection).Count()
	return count + 1
}

// DatabaseInit Initialises the database, and connects to it.
// The collection to use is given by the parameter.
func DatabaseInit(coll string) MongoDB {
	database := MongoDB{
		"ds233763.mlab.com:33763",
		"paragliding_db",
		coll,
		"paragliderAPI",
		"6oLKQOFcxMDCZyd",
	}
	// Connects to the database and returns the struct.
	database.Connect()
	return database
}

// Generates a timestamp for a track.
// The function is monothonic, it wil always count up.
func GenerateTimestamp() int64 {
	// Current time.
	now := time.Now()

	// Unix time in nanoseconds.(Nanoseconds since januar 1970)
	return now.UnixNano()
}
