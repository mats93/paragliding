/*
	File: admin.go
  Contains functions used by API calls to the "Admin paths".
*/

package admin

import (
	"fmt"
	"net/http"

	"github.com/mats93/paragliding/mongodb"
)

// GET: Returns the current count of all tracks in the DB.
// Output: text/plain
func GetTrackCount(w http.ResponseWriter, r *http.Request) {
	database := mongodb.DatabaseInit()

	database.Connect()
	count, _ := database.GetCount()

	fmt.Println("Ant fra admin:", count)
}

// DELETE: Deletes all tracks in the DB.
// Output: text/plain
func DeleteAllTracks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Kommer")
}
