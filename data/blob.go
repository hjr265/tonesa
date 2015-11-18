package data

import (
	"time"
)

type Blob struct {
	Path string `bson:"path"`
	Size int64  `bson:"size"`
}

func (b Blob) SignedURL() string {
	return Bucket.SignedURL(b.Path, time.Now().Add(2*time.Hour))
}
