package ui

import (
	"net/http"
)

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	err := TplIndex.Execute(w, TplIndexValues{})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func init() {
	Router.NewRoute().
		Methods("GET").
		Path("/").
		HandlerFunc(ServeIndex)
}
