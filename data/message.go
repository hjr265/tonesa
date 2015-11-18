package data

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Message struct {
	ID bson.ObjectId `bson:"_id"`

	UploadID bson.ObjectId `bson:"uploadID"`

	AuthorName string `bson:"anonName"`
	Content    string `bson:"content"`

	CreatedAt  time.Time `bson:"createdAt"`
	ModifiedAt time.Time `bson:"modifiedAt"`
}

func listMessagesByQuery(q bson.M, s []string, limit int) ([]Message, error) {
	msgs := []Message{}
	err := sess.DB("").C(messageC).Find(q).Sort(s...).Limit(limit).All(&msgs)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func ListMessagesByUploadID(uplID bson.ObjectId, since time.Time, limit int) ([]Message, error) {
	return listMessagesByQuery(bson.M{"uploadID": uplID, "createdAt": bson.M{"$gt": since}}, []string{"createdAt"}, limit)
}

func ListRecentMessagesByUploadID(uplID bson.ObjectId, limit int) ([]Message, error) {
	return listMessagesByQuery(bson.M{"uploadID": uplID}, []string{"-createdAt"}, limit)
}

func (m *Message) Put() error {
	m.ModifiedAt = time.Now()
	if m.ID == "" {
		m.ID = bson.NewObjectId()
		m.CreatedAt = m.ModifiedAt
	}
	_, err := sess.DB("").C(messageC).UpsertId(m.ID, m)
	return err
}
