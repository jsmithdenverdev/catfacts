package twilio

type Operation = string

const (
	OptOut        Operation = "OPTOUT"
	OptIn         Operation = "OPTIN"
	InvalidAction Operation = "INVALID"
)

func GetTwilioAction(body string) Operation {
	// optOut and optIn taken from
	// https://support.twilio.com/hc/en-us/articles/223134027-Twilio-support-for-opt-out-keywords-SMS-STOP-filtering-
	optOut := []string{
		"STOP",
		"STOPALL",
		"UNSUBSCRIBE",
		"CANCEL",
		"END",
		"QUIT",
	}
	optIn := []string{
		"START",
		"YES",
		"UNSTOP",
	}

	// check for opt out values
	for _, value := range optOut {
		if body == value {
			return OptOut
		}
	}

	// check for opt in values
	for _, value := range optIn {
		if body == value {
			return OptIn
		}
	}

	// if there are no matches return invalid operation
	return InvalidAction
}
