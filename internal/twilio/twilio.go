package twilio

import (
	"fmt"
	"github.com/sfreiberg/gotwilio"
	"gitlab.com/jsmithdenverdev/catfacts/internal/fact"
	"gitlab.com/jsmithdenverdev/catfacts/internal/subscriber"
)

type twilioFactSender struct {
	client *gotwilio.Twilio
	from   string
}

func (t *twilioFactSender) Send(c subscriber.Contact, f fact.Fact) error {
	_, exception, err := t.client.SendSMS(t.from, c, f, "", "")

	if exception != nil {
		return fmt.Errorf("could not send fact with twilio: %w", exception)
	}

	if err != nil {
		return fmt.Errorf("could not send fact with twilio: %w", err)
	}

	return nil
}

func NewFactSender(sid string, token string, from string) fact.Sender {
	client := gotwilio.NewTwilioClient(sid, token)

	return &twilioFactSender{
		client,
		from,
	}
}
