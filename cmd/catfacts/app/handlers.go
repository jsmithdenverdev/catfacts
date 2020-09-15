package app

import (
	"fmt"
	"gitlab.com/jsmithdenverdev/catfacts/internal/subscriber"
	"gitlab.com/jsmithdenverdev/catfacts/internal/twilio"
	"log"
	"net/http"
)

type manageSubscriptionHandler struct {
	service subscriber.SubscriberService
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

	// perform switch logic on the twilio action (OptIn, OptOut, or InvalidOp)
	switch getTwilioAction(operation) {
	case OptIn:
		// create the subscriber
		err = h.service.CreateSubscriber(contact)

		// return an error if creation failed
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// write a twiml response
		writeTwiml(w, "Meow! Welcome to CatFacts! =^._.^=")
	case OptOut:
		// delete the subscriber
		err = h.service.DeleteSubscriber(contact)

		// return an error if deleting failed
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// write a twiml response
		writeTwiml(w, "Meow! Farewell from CatFacts! =^._.^=")
	case InvalidOp:
		// create and return an invalid operation error
		err = fmt.Errorf("request not understood")
		http.Error(w, err.Error(), http.StatusBadRequest)

		// write a twiml response
		writeTwiml(w, "Meow! That request was not understood. =^._.^=")
	}
}

type TwilioOp = string

const (
	OptOut    TwilioOp = "OPTOUT"
	OptIn     TwilioOp = "OPTIN"
	InvalidOp TwilioOp = "INVALID"
)

func getTwilioAction(body string) TwilioOp {
	// optOut and optIn taken from
	// https://support.twilio.com/hc/en-us/articles/223134027-Twilio-support-for-opt-out-keywords-SMS-STOP-filtering-
	optOut := []string{
		"STOP",
		"STOPALL",
		"UNSUBSCRIBE",
		"CANCEL",
		"END",
		"QUIT",
	}
	optIn := []string{
		"START",
		"YES",
		"UNSTOP",
	}

	// check for opt out values
	for _, value := range optOut {
		if body == value {
			return OptOut
		}
	}

	// check for opt in values
	for _, value := range optIn {
		if body == value {
			return OptIn
		}
	}

	// if there are no matches return invalid operation
	return InvalidOp
}

func writeTwiml(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "text/xml")

	response, err := twilio.GenerateTwimlResponse(message)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	_, err = w.Write(response)

	if err != nil {
		log.Fatalf("could not write twiml response: %s", err.Error())
	}
}
