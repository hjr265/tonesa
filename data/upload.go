package data

import (
	"crypto/rand"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	MaxUploadContentSize = 1 * (1 << 20)

	ShortIDAlphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// Upload represents a user uploaded media.
type Upload struct {
	ID      bson.ObjectId `bson:"_id"`
	ShortID string        `bson:"shortID"`

	Kind Kind `bson:"kind"`

	Content Blob `bson:"content"`

	CreatedAt  time.Time `bson:"createdAt"`
	ModifiedAt time.Time `bson:"modifiedAt"`
}

func getUploadByQuery(q bson.M) (*Upload, error) {
	upl := Upload{}
	err := sess.DB("").C(uploadC).Find(q).One(&upl)
	if err == mgo.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &upl, nil
}

// GetUpload finds and returns an upload by its ID.
func GetUpload(id bson.ObjectId) (*Upload, error) {
	return getUploadByQuery(bson.M{"_id": id})
}

// GetUploadByShortID finds and returns an upload by its short ID.
func GetUploadByShortID(shortID string) (*Upload, error) {
	return getUploadByQuery(bson.M{"shortID": shortID})
}

// ResetShortID sets a new random short ID.
func (u *Upload) ResetShortID() error {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return err
	}
	for i := range b {
		b[i] = ShortIDAlphabet[int(b[i])%len(ShortIDAlphabet)]
	}
	u.ShortID = string(b)
	return nil
}

// Put inserts or updates u to MongoDB.
func (u *Upload) Put() error {
	u.ModifiedAt = time.Now()
	if u.ID == "" {
		u.ID = bson.NewObjectId()
		u.CreatedAt = u.ModifiedAt
	}
	if u.ShortID == "" {
		err := u.ResetShortID()
		if err != nil {
			return err
		}
	}
	for {
		_, err := sess.DB("").C(uploadC).UpsertId(u.ID, u)
		if err != nil {
			if mgo.IsDup(err) && strings.Contains(err.Error(), ".shortID_1 dup") {
				err = u.ResetShortID()
				if err != nil {
					return err
				}
				continue
			}
			return err
		}
		break
	}
	return nil
}
