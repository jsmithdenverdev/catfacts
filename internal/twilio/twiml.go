package twilio

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
)

type Response struct {
	Message string `xml:"Message>Body"`
}

func GenerateTwiml(message string) ([]byte, error) {

	response := Response{
		Message: message,
	}

	x, err := xml.MarshalIndent(response, "", "  ")

	if err != nil {
		return nil, fmt.Errorf("could not create twiml response: %w", err)
	}

	return x, nil
}

func WriteTwiml(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "text/xml")

	response, err := GenerateTwiml(message)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	_, err = w.Write(response)

	if err != nil {
		log.Fatalf("could not write twiml response: %s", err.Error())
	}
}
