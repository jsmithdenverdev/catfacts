package rest

import (
	"encoding/xml"
	"net/http"
)

type Response struct {
	Message string `xml:"Message>Body"`
}

func WriteTwimlResponse(w http.ResponseWriter, message string) error {
	x, err := GenerateTwiml(message)

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/xml")

	_, err = w.Write(x)

	if err != nil {
		return err
	}

	return nil
}

func GenerateTwiml(message string) ([]byte, error) {
	response := Response{
		message,
	}

	return xml.MarshalIndent(response, "", "  ")
}
