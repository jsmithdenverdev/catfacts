package internal

import (
	"fmt"
	"github.com/sfreiberg/gotwilio"
)

type Distributor interface {
	Distribute(to string, message string) error
}

type twilioDistributor struct {
	from   string
	client *gotwilio.Twilio
}

func (t twilioDistributor) Distribute(to string, message string) error {
	_, e, err := t.client.SendSMS(t.from, to, message, "", "")

	if e != nil {
		return fmt.Errorf("could not send sms with twilio: %w", e)
	}

	if err != nil {
		return fmt.Errorf("could not send sms with twilio: %w", err)
	}

	return nil
}

func NewTwilioDistributor(sid string, token string, from string) Distributor {
	client := gotwilio.NewTwilioClient(sid, token)

	return twilioDistributor{
		from,
		client,
	}
}
