package twilio

import (
	"github.com/jsmithdenverdev/catfacts"
	"github.com/sfreiberg/gotwilio"
)

type SMSSender struct {
	from   string
	client *gotwilio.Twilio
}

func NewSMSSender(from string, sid string, token string) SMSSender {
	client := gotwilio.NewTwilioClient(sid, token)

	return SMSSender{
		from,
		client,
	}
}

func (sender SMSSender) Send(sms catfacts.SMS) error {
	_, ex, err := sender.client.SendSMS(sender.from, sms.To, sms.Body, "", "")

	if ex != nil {
		return ex
	}

	if err != nil {
		return err
	}

	return nil
}
