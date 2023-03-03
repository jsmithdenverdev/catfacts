package catfacts

import (
	"github.com/chuckpreslar/emission"
	"strings"
)

type SMS struct {
	To   string
	From string
	Body string
}

const (
	SMSBodyWelcome = "😸 Meow! Welcome to CatFacts! You'll be texted a random fact once a day! Reply STOP at any time to unsubscribe."
	SMSBodyGoodbye = "😿 Farewell for now! You can reply START at any time to re-subscribe."
	SMSBodyHelp    = "🙀 Does someone need some help?! This is CatFacts. Reply START to start receiving a random cat fact once a day. Already subscribed? Reply STOP at any time to unsubscribe."
	SMSBodyCuss    = "😾 $#*!"
	SMSBodyError   = "😾 $#*! Something went wrong! Server Cat is working on fixing it now (and boy is he pissed)."
	SMSBodyFact    = "😺 Here's your daily fact. "
)

func HandleSMS(emitter *emission.Emitter) func(request SMS) {
	return func(request SMS) {
		var requestMap = map[any][]string{
			CommandSubscribe:   {"start", "subscribe", "unstop"},
			CommandUnsubscribe: {"stop", "stopall", "unsubscribe", "end", "quit"},
			CommandCuss:        {"cuss"},
			QueryHelp:          {"help", "info"},
		}

		check := strings.ToLower(request.Body)
		check = strings.TrimSpace(check)

		for match, keywords := range requestMap {
			for _, keyword := range keywords {
				if check == keyword {
					emitter.Emit(match, request)
				}
			}
		}
	}
}
