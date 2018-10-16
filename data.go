/*
  File: data.go
  Contains constant variables, structs and global variables for the API.
*/

package main

import "time"

// INFORMATION = The API inforomation text.
const INFORMATION = "Service for Paragliding tracks"

// VERSION = The API version that is currently running.
const VERSION = "v1"

// Format for the IGC file.
type track struct {
	ID          int       `json:"-"`
	HDate       time.Time `json:"H_date"`
	Pilot       string    `json:"pilot"`
	Glider      string    `json:"glider"`
	GliderID    string    `json:"glider_id"`
	TrackLength float64   `json:"track_length"`
	TrackSrcURL string    `json:"track_src_url"`
}

// Format for the API information.
type apiInfo struct {
	Uptime  string `json:"uptime"`
	Info    string `json:"info"`
	Version string `json:"version"`
}

// Format for the url information.
type url struct {
	URL string `json:"url"`
}

// Format for the id to be returned when a new track is posted.
type id struct {
	ID int `json:"id"`
}

// Slice with structs of type "track" (in-memory storage).
var trackSlice []track

// Generates an unique ID for all tracks.
var lastUsedID int

// The start time for the API service.
var startTime time.Time = time.Now()
