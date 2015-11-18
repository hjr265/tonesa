package ui

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hjr265/tonesa/data"
)

func ServeUpload(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	upl, err := data.GetUploadByShortID(vars["shortID"])
	if err != nil {
		ServeInternalServerError(w, r)
		return
	}
	if upl == nil {
		ServeNotFound(w, r)
		return
	}

	err = TplUploadView.Execute(w, TplUploadViewValues{
		Upload: upl,
	})
	if err != nil {
		ServeInternalServerError(w, r)
		return
	}
}

func init() {
	Router.NewRoute().
		Methods("GET").
		Path("/u/{shortID}").
		HandlerFunc(ServeUpload)
}
