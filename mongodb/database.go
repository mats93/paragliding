/*
	File: database.go
  Handles the mongoDB.
*/

package mongodb

import (
	"fmt"

	"github.com/globalsign/mgo"
	//"github.com/globalsign/mgo/bson"
)

// Test something.
func Test() {
	session, _ := mgo.Dial("mongodb://admin:passord1@ds233763.mlab.com:33763/paragliding_db")
	fmt.Println(session)
}
