package twilio

import "encoding/xml"

type Response struct {
	Message string `xml:"Message>Body"`
}

func GenerateTwiml(message string) ([]byte, error) {
	response := Response{
		message,
	}

	return xml.MarshalIndent(response, "", "  ")
}
