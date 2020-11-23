package twilio

import "github.com/sfreiberg/gotwilio"

type SmsSender struct {
	client *gotwilio.Twilio
	from   string
}

func NewSmsSender(sid, token, from string) SmsSender {
	client := gotwilio.NewTwilioClient(sid, token)
	return SmsSender{
		client,
		from,
	}
}

func (s SmsSender) Send(to, message string) error {
	_, ex, err := s.client.SendSMS(s.from, to, message, "", "")
	if ex != nil {
		// this error is returned if a message is accidentally sent to an unsubscribed party
		if ex.Code == 21211 {
			return nil
		}
		return ex
	}

	if err != nil {
		return err
	}

	return nil
}
