package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gitlab.com/jsmithdenverdev/catfacts/internal/app"
)

func CreateRouter(a app.App) *mux.Router {
	r := mux.NewRouter()

	_ = handlers.CORS()

	r.HandleFunc("/subscribers", a.ManageSubscription).Methods("POST")
	r.HandleFunc("/facts/distribute", a.SendFactToSubscribers).Methods("POST")

	return r
}
