package rest

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

type Response struct {
	Message string `xml:"Message>Body"`
}

func WriteTwimlResponse(w http.ResponseWriter, message string) error {
	response := Response{
		message,
	}

	x, err := xml.MarshalIndent(response, "", "  ")

	if err != nil {
		return fmt.Errorf("could not create twiml response: %w", err)
	}

	w.Header().Set("Content-Type", "text/xml")

	_, err = w.Write(x)

	if err != nil {
		return err
	}

	return nil
}
