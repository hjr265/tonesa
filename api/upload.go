package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/AdRoll/goamz/s3"
	"github.com/hjr265/tonesa/data"
)

type Upload struct {
	ID      string        `json:"id"`
	ShortID string        `json:"shortID"`
	Content UploadContent `json:"content"`
}

type UploadContent struct {
	URL string `json:"url"`
}

func HandleUploadCreate(w http.ResponseWriter, r *http.Request) {
	f, h, err := r.FormFile("file")
	if err != nil {
		ServeBadRequest(w, r)
		return
	}

	b := bytes.Buffer{}
	n, err := io.Copy(&b, io.LimitReader(f, data.MaxUploadContentSize+10))
	if err != nil {
		ServeInternalServerError(w, r)
		return
	}
	if n > data.MaxUploadContentSize {
		ServeBadRequest(w, r)
		return
	}

	id := bson.NewObjectId()
	upl := data.Upload{
		ID:   id,
		Kind: data.Image,
		Content: data.Blob{
			Path: "/uploads/" + id.Hex(),
			Size: n,
		},
	}

	err = data.Bucket.Put(upl.Content.Path, b.Bytes(), h.Header.Get("Content-Type"), s3.Private, s3.Options{})
	if err != nil {
		ServeInternalServerError(w, r)
		return
	}

	err = upl.Put()
	if err != nil {
		ServeInternalServerError(w, r)
		return
	}

	res := Upload{
		ID:      upl.ID.Hex(),
		ShortID: upl.ShortID,
		Content: UploadContent{
			URL: upl.Content.SignedURL(),
		},
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		ServeInternalServerError(w, r)
		return
	}
}

func init() {
	Router.NewRoute().
		Methods("POST").
		Path("/uploads").
		HandlerFunc(HandleUploadCreate)
}
