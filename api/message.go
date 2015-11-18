package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
	"github.com/hjr265/tonesa/data"
	"github.com/hjr265/tonesa/hub"
)

type Message struct {
	ID         string `json:"id"`
	AuthorName string `json:"authorName"`
	Content    string `json:"content"`
	CreatedAt  string `json:"createdAt"`
}

func ServeMessageList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	if !bson.IsObjectIdHex(idStr) {
		ServeNotFound(w, r)
		return
	}
	upl, err := data.GetUpload(bson.ObjectIdHex(idStr))
	if err != nil {
		ServeInternalServerError(w, r)
		return
	}
	if upl == nil {
		ServeNotFound(w, r)
		return
	}

	sinceStr := r.URL.Query().Get("since")

	var msgs []data.Message
	if sinceStr != "" {
		since, err := time.Parse(time.RFC3339, sinceStr)
		if err != nil {
			ServeBadRequest(w, r)
			return
		}

		msgs, err = data.ListMessagesByUploadID(upl.ID, since, 16)
		if err != nil {
			ServeInternalServerError(w, r)
			return
		}

	} else {
		msgs, err = data.ListRecentMessagesByUploadID(upl.ID, 16)
		if err != nil {
			ServeInternalServerError(w, r)
			return
		}
	}

	col := []Message{}
	for _, msg := range msgs {
		res := Message{
			ID:         msg.ID.Hex(),
			AuthorName: msg.AuthorName,
			Content:    msg.Content,
			CreatedAt:  msg.CreatedAt.Format(time.RFC3339),
		}
		col = append(col, res)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(col)
	if err != nil {
		ServeInternalServerError(w, r)
		return
	}
}

func HandleMessageCreate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	if !bson.IsObjectIdHex(idStr) {
		ServeNotFound(w, r)
		return
	}
	upl, err := data.GetUpload(bson.ObjectIdHex(idStr))
	if err != nil {
		ServeInternalServerError(w, r)
		return
	}
	if upl == nil {
		ServeNotFound(w, r)
		return
	}

	body := Message{}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		ServeBadRequest(w, r)
		return
	}

	msg := data.Message{}
	msg.UploadID = upl.ID
	msg.AuthorName = body.AuthorName
	msg.Content = body.Content
	msg.Put()

	res := Message{
		ID:         msg.ID.Hex(),
		AuthorName: msg.AuthorName,
		Content:    msg.Content,
		CreatedAt:  msg.CreatedAt.Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		ServeInternalServerError(w, r)
		return
	}

	err = hub.Emit("upload:"+upl.ID.Hex(), "message:"+msg.ID.Hex())
	if err != nil {
		log.Print(err)
	}
}

func init() {
	Router.NewRoute().
		Methods("GET").
		Path("/uploads/{id}/messages").
		HandlerFunc(ServeMessageList)
	Router.NewRoute().
		Methods("POST").
		Path("/uploads/{id}/messages").
		HandlerFunc(HandleMessageCreate)
}
