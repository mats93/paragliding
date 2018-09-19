/*
  File: functions_test.go
  Contains unit tests for functions.go
*/

package main

import (
	"testing"
	"time"
)

// Test if the correct ID is returned.
func Test_retriveTrackById(t *testing.T) {

	// Test if it returns an error if the trackSlice is emtpy.
	if _, err := retriveTrackById(10); err == nil {
		t.Error("Error: Did not return error when track array was empty")
	}

	// Adds a track to the trackSlice.
	trackSlice = append(trackSlice,
		track{10, time.Now(), "pilot", "glider", "glider_id", 20.4})

	// Test if it returns an error if the track was not found.
	if _, err := retriveTrackById(1); err == nil {
		t.Error("Error: Did not return error when requestet ID did not excist")
	}

	// Test if the correct track with the specified ID is sent back.
	newTrack, _ := retriveTrackById(10)
	if newTrack.Id != 10 {
		t.Error("Error: Did not return the correct track")
	}
}
