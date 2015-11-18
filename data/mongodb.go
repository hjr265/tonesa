package data

import (
	"gopkg.in/mgo.v2"
)

const (
	uploadC  = "uploads"
	messageC = "messages"
)

var sess *mgo.Session

// OpenSession connects to MongoDB via the given URL.
func OpenSession(url string) (err error) {
	sess, err = mgo.Dial(url)
	return err
}

// MakeIndexes ensures necessary indexes in MongoDB.
func MakeIndexes() error {
	return nil
}
