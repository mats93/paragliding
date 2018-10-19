/*
  File: database_test.go
  Contains unit tests for database.go
*/

package mongodb

import (
	"testing"
	"time"
)

// The collection to test against.
const COLLECTION = "testing"

// Method to test: Insert().
// Test if the correct track is insertet into the database.
func Test_Insert(t *testing.T) {
	// Connects to the database.
	database := DatabaseInit(COLLECTION)

	expected := Track{100, time.Now(), "pilot1", "glider1", "glider_id1", 20.1, "http://test1.test"}

	// Check to see if insert generates error.
	err := database.Insert(expected)
	if err != nil {
		t.Errorf("Method returned an unexpected error: %v", err)
	}

	// Check if the same track was returned.
	actual, _ := database.FindByID(100)
	if actual[0].ID != expected.ID {
		t.Errorf("Method returned wrong Track: got %v want %v",
			actual[0].ID, expected.ID)
	}

	// Deletes all from the database.
	database.DeleteAllTracks()

	// Closes the database session.
	defer MDB.Session.Close()
}

// Method to test: DeleteAllTracks().
// Test if all inserted tracks was deleted from the database.
func Test_DeleteAllTracks(t *testing.T) {
	// Connects to the database.
	database := DatabaseInit(COLLECTION)

	// Expected count when all 5 tracks are deleted.
	expected := 0

	// Inserts 5 tracks to the databae.
	database.Insert(Track{1, time.Now(), "pilot1", "glider1", "glider_id1", 20.1, "http://test1.test"})
	database.Insert(Track{2, time.Now(), "pilot2", "glider2", "glider_id2", 20.2, "http://test2.test"})
	database.Insert(Track{3, time.Now(), "pilot3", "glider3", "glider_id3", 20.3, "http://test3.test"})
	database.Insert(Track{4, time.Now(), "pilot4", "glider4", "glider_id4", 20.4, "http://test4.test"})
	database.Insert(Track{5, time.Now(), "pilot5", "glider5", "glider_id5", 20.5, "http://test5.test"})

	// Check for errors when deleting.
	err := database.DeleteAllTracks()
	if err != nil {
		t.Errorf("Method returned unexpected error: %v", err)
	}

	// Check the count of the database.
	actual, _ := database.GetCount()
	if actual != expected {
		t.Errorf("Method returned wrong count: got %d want %d",
			actual, expected)
	}
}

// Method to test: FindAll().
// Test if an empty Track is returned when collection is empty.
func Test_FindAll_Empty(t *testing.T) {
	// Connects to the database.
	database := DatabaseInit(COLLECTION)

	// Check if the correct track slice is returned (nil).
	actual, _ := database.FindAll()

	if actual != nil {
		t.Errorf("Method did not return empty track: got %v want %s",
			actual, "[]")
	}

	// Closes the database session.
	defer MDB.Session.Close()
}

// Method to test: FindAll().
// Test if the correct tracks are returned from the database.
func Test_FindAll(t *testing.T) {
	// Connects to the database.
	database := DatabaseInit(COLLECTION)

	// Expected results from the database.
	var expected []Track
	expected = append(expected,
		Track{1, time.Now(), "pilot1", "glider1", "glider_id1", 20.1, "http://test1.test"},
		Track{2, time.Now(), "pilot2", "glider2", "glider_id2", 20.2, "http://test2.test"})

	// Test if correct track slice is retunred when collection has data.
	database.Insert(Track{1, time.Now(), "pilot1", "glider1", "glider_id1", 20.1, "http://test1.test"})
	database.Insert(Track{2, time.Now(), "pilot2", "glider2", "glider_id2", 20.2, "http://test2.test"})
	actual, _ := database.FindAll()

	// Check if method did not return an emtpy slice.
	if actual == nil {
		t.Errorf("Method returned empty track slice: \ngot %v \nwant %v",
			actual, expected)
	} else {
		// Check if the slice length is as expected.
		if len(actual) != len(expected) {
			t.Errorf("Track slice is not the same length as expected: got %d want %d",
				len(actual), len(expected))
		} else {
			// Check if the content of the slice is as expected (Not comparing timestamp),
			// therefore the many if statements.
			for i := 0; i < len(expected); i++ {
				if expected[i].ID != actual[i].ID {
					t.Errorf("Method returned wrong ID: got %v want %v",
						actual[i].ID, expected[i].ID)
				}
				if expected[i].Pilot != actual[i].Pilot {
					t.Errorf("Method returned wrong Pilot: got %v want %v",
						actual[i].Pilot, expected[i].Pilot)
				}
				if expected[i].Glider != actual[i].Glider {
					t.Errorf("Method returned wrong Glider: got %v want %v",
						actual[i].Glider, expected[i].Glider)
				}
				if expected[i].GliderID != actual[i].GliderID {
					t.Errorf("Method returned wrong GliderID: got %v want %v",
						actual[i].GliderID, expected[i].GliderID)
				}
				if expected[i].TrackLength != actual[i].TrackLength {
					t.Errorf("Method returned wrong TrackLength: got %v want %v",
						actual[i].TrackLength, expected[i].TrackLength)
				}
				if expected[i].TrackSrcURL != actual[i].TrackSrcURL {
					t.Errorf("Method returned wrong TrackSrcURL: got %v want %v",
						actual[i].TrackSrcURL, expected[i].TrackSrcURL)
				}
			}
		}
	}

	// Deletes all from the database.
	database.DeleteAllTracks()

	// Closes the database session.
	defer MDB.Session.Close()
}

// Method to test: FindByID().
// Test the error message when the collection is empty.
func Test_FindByID_Empty(t *testing.T) {
	// Connects to the database.
	database := DatabaseInit(COLLECTION)

	_, err := database.FindByID(1)
	if err == nil {
		t.Errorf("Method returned wrong error: got %v want %v",
			"not found", err)
	}

	// Closes the database session.
	defer MDB.Session.Close()
}

// Method to test: FindByID().
// Test if the correct track is returned from the database.
func Test_FindByID(t *testing.T) {
	// Connects to the database.
	database := DatabaseInit(COLLECTION)

	// Expected results from the database.
	var expected []Track
	expected = append(expected,
		Track{1, time.Now(), "pilot1", "glider1", "glider_id1", 20.1, "http://test1.test"})

	// Test data.
	database.Insert(Track{1, time.Now(), "pilot1", "glider1", "glider_id1", 20.1, "http://test1.test"})

	// Check if correct track was returned, when given the ID 1.
	actual, _ := database.FindByID(1)

	if actual != nil && (actual[0].ID != expected[0].ID) {
		t.Errorf("Method did not return the correct track ID: got %v want %v",
			actual[0].ID, expected[0].ID)
	}

	// Deletes all from the database.
	database.DeleteAllTracks()

	// Closes the database session.
	defer MDB.Session.Close()
}

// Method to test: GetCount().
// Test if the correct count is returned when the database is empty.
func Test_GetCount_Empty(t *testing.T) {
	// Connects to the database.
	database := DatabaseInit(COLLECTION)

	expected := 0

	actual, err := database.GetCount()
	if err != nil {
		t.Errorf("Method returned an unexpected error: %v", err)
	}

	// Check if correct count was returned (0).
	if actual != expected {
		t.Errorf("Method returned wrong count: got %d want %d",
			actual, expected)
	}

	// Closes the database session.
	defer MDB.Session.Close()
}

// Method to test: GetCount().
// Test if the correct count is returned when there are tracks in the database.
func Test_GetCount(t *testing.T) {
	// Connects to the database.
	database := DatabaseInit(COLLECTION)

	expected := 1

	// Test data.
	database.Insert(Track{1, time.Now(), "pilot1", "glider1", "glider_id1", 20.1, "http://test1.test"})

	actual, err := database.GetCount()
	if err != nil {
		t.Errorf("Method returned an unexpected error: %v", err)
	}

	// Check if correct count was returned (1).
	if actual != expected {
		t.Errorf("Method returned wrong count: got %d want %d",
			actual, expected)
	}

	// Deletes all from the database.
	database.DeleteAllTracks()

	// Closes the database session.
	defer MDB.Session.Close()
}
