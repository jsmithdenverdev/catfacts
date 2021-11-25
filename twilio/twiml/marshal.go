package twiml

import "encoding/xml"

func Marshal(body string) ([]byte, error) {
	twiml := struct {
		XMLName xml.Name `xml:"Response"`
		Body    string   `xml:"Message>Body"`
	}{
		Body: body,
	}

	return xml.MarshalIndent(twiml, " ", " ")
}
