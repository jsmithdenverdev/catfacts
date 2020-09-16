package main

import (
	"fmt"
	"github.com/jsmithdenverdev/catfacts/internal/subscriber"
	"github.com/jsmithdenverdev/catfacts/internal/twilio"
	"net/http"
)

type manageSubscriptionHandler struct {
	service subscriber.Service
}

func (h manageSubscriptionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	// parse the post form
	err = r.ParseForm()

	// returning an error if the form could not be parsed
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// get the twilio operation and the requesting contact from the post form
	operation := r.PostForm.Get("Body")
	contact := r.PostForm.Get("From")

	// if no contact was supplied this is an invalid request
	if len(contact) == 0 {
		http.Error(w, "could not create subscriber: malformed twilio request: missing parameter From", http.StatusBadRequest)
	}

	// perform switch logic on the twilio action (OptIn, OptOut, or InvalidAction)
	switch twilio.GetTwilioAction(operation) {
	case twilio.OptIn:
		// create the subscriber
		err = h.service.CreateSubscriber(contact)

		// return an error if creation failed
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// write a twiml response
		twilio.WriteTwiml(w, "Meow! Welcome to CatFacts! =^._.^=. Cancel your subscription at any time by replying STOP.")
	case twilio.OptOut:
		// delete the subscriber
		err = h.service.DeleteSubscriber(contact)

		// return an error if deleting failed
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// write a twiml response
		twilio.WriteTwiml(w, "Meow! Farewell from CatFacts! =^._.^=")
	case twilio.InvalidAction:
		// create and return an invalid operation error
		err = fmt.Errorf("request not understood")
		http.Error(w, err.Error(), http.StatusBadRequest)

		// write a twiml response
		twilio.WriteTwiml(w, "Meow! That request was not understood. =^._.^=")
	}
}
