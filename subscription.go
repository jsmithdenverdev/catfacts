package catfacts

import "strings"

type action = int

const (
	actionSubscribe   action = iota
	actionUnsubscribe action = iota
	actionHelp        action = iota
	actionUnknown     action = iota
)

type Reply = string

const (
	ReplyWelcome       Reply = "😸 Meow! Welcome to cat facts. You can unsubscribe at any time by replying stop."
	ReplyGoodbye             = "😿 Farewell from cat facts! You can re-subscribe at any time by replying start."
	ReplyHelp                = "🐾 This is cat facts. Reply start to start receiving a cat fact by text once daily. Reply stop to cancel your subscription."
	ReplyUnknown             = "🙀 Server cat could not understand your request! Reply help for more information."
	ReplyAlreadyExists       = "🙀 A subscription was already found for this number!"
)

type Subscription struct {
	Contact string `json:"contact"`
}

type SubscriptionService interface {
	CreateSubscription(subscription Subscription) error
	DeleteSubscription(contact string) error
	All() ([]*Subscription, error)
}

func ManageSMSSubscription(s SubscriptionService, message, contact string) (Reply, error) {
	switch getSubscriptionAction(message) {
	case actionSubscribe:
		subscription := Subscription{
			Contact: contact,
		}
		if err := s.CreateSubscription(subscription); err != nil {
			if err == ErrSubscriptionExists {
				return ReplyAlreadyExists, ErrSubscriptionExists
			}
			return "", err
		}
		return ReplyWelcome, nil
	case actionUnsubscribe:
		if err := s.DeleteSubscription(contact); err != nil {
			return "", err
		}
		return ReplyGoodbye, nil
	case actionHelp:
		return ReplyHelp, nil
	case actionUnknown:
		return ReplyUnknown, nil
	}

	return ReplyUnknown, nil
}

func getSubscriptionAction(input string) action {
	actionMap := map[action][]string{
		actionSubscribe:   {"START", "YES", "UNSTOP"},
		actionUnsubscribe: {"STOP", "STOPALL", "UNSUBSCRIBE", "CANCEL", "END", "QUIT"},
		actionHelp:        {"HELP", "INFO"},
	}

	for action, keywords := range actionMap {
		for _, keyword := range keywords {
			if keyword == strings.ToUpper(input) {
				return action
			}
		}
	}

	return actionUnknown
}
