/*
  File: data.go
  Contains constant variables, structs and global variables for the API.
*/

package main

import "time"

// Const for the API information.
const INFORMATION = "Service for IGC tracks"
const VERSION = "v1"

// Format for the IGC file.
type track struct {
	Id           int       `json:"-"`
	H_date       time.Time `json:"H_date"`
	Pilot        string    `json:"pilot"`
	Glider       string    `json:"glider"`
	Glider_id    string    `json:"glider_id"`
	Track_length float64   `json:"track_length"`
}

// Format for the API information.
type apiInfo struct {
	Uptime  string `json:"uptime"`
	Info    string `json:"info"`
	Version string `json:"version"`
}

// Format for the url information.
type url struct {
	Url string `json:"url"`
}

// Format for the id to be returned when a new track is posted.
type id struct {
	Id int `json:"id"`
}

// Slice with structs of type "track" (in-memory storage).
var trackSlice []track

// Generates an unique ID for all tracks.
var lastUsedID int = 0

// The start time for the API service.
var startTime time.Time = time.Now()
