package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jsmithdenverdev/catfacts/pkg/catfacts"
)

type SubscriptionHandler struct {
	Service   catfacts.SubscriptionStore
	SMSSender catfacts.SMSSender
}

func (handler SubscriptionHandler) HandleCreate(rw http.ResponseWriter, req *http.Request) {
	var request struct {
		Contact string `json:"contact"`
	}

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	err := handler.Service.Insert(catfacts.Subscription{
		Contact: request.Contact,
	})

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	if err := handler.SMSSender.Send(catfacts.SMS{
		To:   request.Contact,
		Body: catfacts.ReplyWelcome,
	}); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.WriteHeader(http.StatusCreated)
}

func (handler SubscriptionHandler) HandleDelete(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	contact := vars["contact"]

	if err := handler.Service.Delete(contact); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.WriteHeader(http.StatusOK)
}

func (handler SubscriptionHandler) HandleList(rw http.ResponseWriter, req *http.Request) {
	var response struct {
		Subscriptions []catfacts.Subscription `json:"subscriptions"`
	}

	subscriptions, err := handler.Service.All()

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	response.Subscriptions = subscriptions

	rw.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(rw).Encode(response); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}
