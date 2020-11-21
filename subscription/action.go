package subscription

import "strings"

type Action = int

const (
	ActionSubscribe = iota
	ActionUnsubscribe
	ActionHelp
	ActionUnknown
)

func GetAction(message string) Action {
	actionMap := map[Action][]string{
		ActionSubscribe:   {"START", "YES", "UNSTOP"},
		ActionUnsubscribe: {"STOP", "STOPALL", "UNSUBSCRIBE", "CANCEL", "END", "QUIT"},
		ActionHelp:        {"HELP", "INFO"},
	}

	for action, keywords := range actionMap {
		for _, keyword := range keywords {
			if keyword == strings.ToUpper(message) {
				return action
			}
		}
	}

	return ActionUnknown
}
