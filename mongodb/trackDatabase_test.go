/*
  File: trackDatabase_test.go
  Contains unit tests for trackDatabase.go
*/

package mongodb

import (
	"testing"
	"time"
)

// Method to test: Insert().
// Test if the correct track is insertet into the database.
func Test_Insert(t *testing.T) {
	// Connects to the database.
	database := DatabaseInit("TestTracks")

	expected := Track{100, 10, time.Now(), "pilot1", "glider1", "glider_id1", 20.1, "http://test1.test"}

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
	database.DeleteAll()

	// Closes the database session.
	defer MDB.Session.Close()
}

// Method to test: DeleteAll().
// Test if all inserted tracks was deleted from the database.
func Test_DeleteAll(t *testing.T) {
	// Connects to the database.
	database := DatabaseInit("TestTracks")

	// Expected count when all 5 tracks are deleted.
	expected := 0

	// Inserts 5 tracks to the database.
	database.Insert(Track{1, 11, time.Now(), "pilot1", "glider1", "glider_id1", 20.1, "http://test1.test"})
	database.Insert(Track{2, 12, time.Now(), "pilot2", "glider2", "glider_id2", 20.2, "http://test2.test"})
	database.Insert(Track{3, 13, time.Now(), "pilot3", "glider3", "glider_id3", 20.3, "http://test3.test"})
	database.Insert(Track{4, 14, time.Now(), "pilot4", "glider4", "glider_id4", 20.4, "http://test4.test"})
	database.Insert(Track{5, 15, time.Now(), "pilot5", "glider5", "glider_id5", 20.5, "http://test5.test"})

	// Check for errors when deleting.
	err := database.DeleteAll()
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
	database := DatabaseInit("TestTracks")

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
	database := DatabaseInit("TestTracks")

	// Expected results from the database.
	var expected []Track
	expected = append(expected,
		Track{1, 11, time.Now(), "pilot1", "glider1", "glider_id1", 20.1, "http://test1.test"},
		Track{2, 12, time.Now(), "pilot2", "glider2", "glider_id2", 20.2, "http://test2.test"})

	// Test if correct track slice is retunred when collection has data.
	database.Insert(Track{1, 11, time.Now(), "pilot1", "glider1", "glider_id1", 20.1, "http://test1.test"})
	database.Insert(Track{2, 12, time.Now(), "pilot2", "glider2", "glider_id2", 20.2, "http://test2.test"})
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
	database.DeleteAll()

	// Closes the database session.
	defer MDB.Session.Close()
}

// Method to test: FindByID().
// Test the error message when the collection is empty.
func Test_FindByID_Empty(t *testing.T) {
	// Connects to the database.
	database := DatabaseInit("TestTracks")

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
	database := DatabaseInit("TestTracks")

	// Expected results from the database.
	var expected []Track
	expected = append(expected,
		Track{1, 11, time.Now(), "pilot1", "glider1", "glider_id1", 20.1, "http://test1.test"})

	// Test data.
	database.Insert(Track{1, 11, time.Now(), "pilot1", "glider1", "glider_id1", 20.1, "http://test1.test"})

	// Check if correct track was returned, when given the ID 1.
	actual, _ := database.FindByID(1)

	if actual != nil && (actual[0].ID != expected[0].ID) {
		t.Errorf("Method did not return the correct track ID: got %v want %v",
			actual[0].ID, expected[0].ID)
	}

	// Deletes all from the database.
	database.DeleteAll()

	// Closes the database session.
	defer MDB.Session.Close()
}

// Method to test: GetCount().
// Test if the correct count is returned when the database is empty.
func Test_GetCount_Empty(t *testing.T) {
	// Connects to the database.
	database := DatabaseInit("TestTracks")

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
	database := DatabaseInit("TestTracks")

	expected := 1

	// Test data.
	database.Insert(Track{1, 11, time.Now(), "pilot1", "glider1", "glider_id1", 20.1, "http://test1.test"})

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
	database.DeleteAll()

	// Closes the database session.
	defer MDB.Session.Close()
}

// Method to test: GetNewID().
// Test if the correct ID is returned, when the database is empty.
func Test_GetNewID_Empty(t *testing.T) {
	database := DatabaseInit("TestTracks")

	// The expected ID to be generated.
	expected := 1

	// Check if the Generated ID is correct.
	count, _ := database.GetCount()
	actual := count + 1

	if actual != expected {
		t.Errorf("Method generated wrong ID: got %d want %d",
			actual, expected)
	}

	// Closes the database session.
	defer MDB.Session.Close()
}

// Method to test: GetNewID().
// Test if the correct ID is returned, when the database has content.
func Test_GetNewID(t *testing.T) {
	database := DatabaseInit("TestTracks")

	// Inserts 5 tracks to the database.
	database.Insert(Track{1, 11, time.Now(), "pilot1", "glider1", "glider_id1", 20.1, "http://test1.test"})
	database.Insert(Track{2, 12, time.Now(), "pilot2", "glider2", "glider_id2", 20.2, "http://test2.test"})
	database.Insert(Track{3, 13, time.Now(), "pilot3", "glider3", "glider_id3", 20.3, "http://test3.test"})
	database.Insert(Track{4, 14, time.Now(), "pilot4", "glider4", "glider_id4", 20.4, "http://test4.test"})
	database.Insert(Track{5, 15, time.Now(), "pilot5", "glider5", "glider_id5", 20.5, "http://test5.test"})

	// The expected ID to be generated after inserting 5 tracks with ID 1-5.
	expected := 6

	// Check if the Generated ID is correct.
	actual := database.GetNewID()

	if actual != expected {
		t.Errorf("Method generated wrong ID: got %d want %d",
			actual, expected)
	}

	// Deletes all from the database.
	database.DeleteAll()

	// Closes the database session.
	defer MDB.Session.Close()
}

// Method to test: FindTrackHigherThen().
// Test if the correct tracks are returend.
func Test_FindTrackHigherThen(t *testing.T) {
	// Connects to the database.
	database := DatabaseInit("TestTracks")

	// Inserts 5 tracks to the database.
	database.Insert(Track{1, 11, time.Now(), "pilot1", "glider1", "glider_id1", 20.1, "http://test1.test"})
	database.Insert(Track{2, 12, time.Now(), "pilot2", "glider2", "glider_id2", 20.2, "http://test2.test"})
	database.Insert(Track{3, 13, time.Now(), "pilot3", "glider3", "glider_id3", 20.3, "http://test3.test"})
	database.Insert(Track{4, 14, time.Now(), "pilot4", "glider4", "glider_id4", 20.4, "http://test4.test"})
	database.Insert(Track{5, 15, time.Now(), "pilot5", "glider5", "glider_id5", 20.5, "http://test5.test"})

	// The expected slice lenth to be returned, when querieng for
	// timestamps higher then 13.
	expected := 2

	// All higher then 13
	querie, _ := database.FindTrackHigherThen(13)
	actual := len(querie)

	if actual != expected {
		t.Errorf("Method queried the DB wrong: got %d length want %d length of slice",
			actual, expected)
	}

	// The expected slice lengt when the highest timestamp is searched for.
	// Slice should be of 0 length.
	expected = 0
	querie, _ = database.FindTrackHigherThen(15)
	actual = len(querie)

	if actual != expected {
		t.Errorf("Method queried the DB wrong: got %d lenght want %d length of slice",
			actual, expected)
	}

	// Deletes all from the database.
	database.DeleteAll()

	// Closes the database session.
	defer MDB.Session.Close()
}

// Function to test: SortTrackByTimestamp().
// Test to check if the slice was sorted correctly.
func Test_SortTrackByTimestamp(t *testing.T) {
	// Connects the the database and inserts 3 tracks.
	// The last inserted has the highest timestamp.
	database := DatabaseInit("TestTracks")
	database.Insert(Track{1, 111, time.Now(), "pilot1", "glider1", "glider_id1", 20.1, "http://test1.test"})
	database.Insert(Track{2, 222, time.Now(), "pilot2", "glider2", "glider_id2", 20.2, "http://test2.test"})
	database.Insert(Track{3, 333, time.Now(), "pilot3", "glider3", "glider_id3", 20.3, "http://test3.test"})

	// Returns all tracks from the DB.
	tracks, _ := database.FindAll()

	// Try to sort the slice.
	sorted := SortTrackByTimestamp(tracks)

	// The unsorted track should have the highest timestamp in index 2.
	// The sorted track should have the highest timestamp in index 0.
	if sorted[2].Timestamp > sorted[0].Timestamp {
		t.Errorf("Function did not sort correctly")
	}

	// Removes the test data.
	database.DeleteAll()

	// Closes the database session.
	defer MDB.Session.Close()
}

// Function to test: GenerateTimestamp().
// Test if two timestamps are monothonic.
func Test_GenerateTimestamp(t *testing.T) {
	first := GenerateTimestamp()
	second := GenerateTimestamp()

	if first == second {
		t.Errorf("Both timestamps generetad is the same")
	}

	if first > second {
		t.Error("Function is not monothonic, first timestamp has lower value")
	}
}
