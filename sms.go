package catfacts

import "strings"

type SMS struct {
	To   string
	From string
	Body string
}

type SMSSender interface {
	Send(sms SMS) error
}

type SMSCommand int

const (
	CommandSubscribe SMSCommand = iota
	CommandUnsubscribe
	CommandHelp
	CommandCuss
)

const (
	SMSBodyWelcome = "😸 Meow! Welcome to CatFacts! You'll be texted a random fact once a day! Reply STOP at any time to unsubscribe."
	SMSBodyGoodbye = "😿 Farewell for now! You can reply START at any time to re-subscribe."
	SMSBodyHelp    = "🙀 Does someone need some help?! This is CatFacts. Reply START to start receiving a random cat fact once a day. Already subscribed? Reply STOP at any time to unsubscribe."
	SMSBodyCuss    = "😾 $#*!"
	SMSBodyError   = "😾 $#*! Something went wrong! Server Cat is working on fixing it now (and boy is he pissed)."
	SMSBodyFact    = "😺 Here's your daily fact. "
)

func (sms SMS) ToCommand() SMSCommand {
	var commandMap map[SMSCommand][]string = map[SMSCommand][]string{
		CommandSubscribe:   {"start", "subscribe", "unstop"},
		CommandUnsubscribe: {"stop", "stopall", "unsubscribe", "end", "quit"},
		CommandHelp:        {"help", "info"},
		CommandCuss:        {"cuss"},
	}

	check := strings.ToLower(sms.Body)
	check = strings.TrimSpace(check)

	for command, keywords := range commandMap {

		for _, keyword := range keywords {
			if check == keyword {
				return command
			}
		}
	}

	return CommandHelp
}
