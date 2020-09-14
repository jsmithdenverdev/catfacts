package internal

import (
	"encoding/xml"
	"fmt"
)

type Response struct {
	Message string `xml:"Message>Body"`
}

func GenerateTwimlResponse(message string) ([]byte, error) {

	response := Response{
		Message: message,
	}

	x, err := xml.MarshalIndent(response, "", "  ")

	if err != nil {
		return nil, fmt.Errorf("could not create twiml response: %w", err)
	}

	return x, nil
}
