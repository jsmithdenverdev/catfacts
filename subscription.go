package catfacts

import (
	"errors"
	"github.com/chuckpreslar/emission"
)

type Subscription struct {
	Contact string
}

// validate performs basic validation on a Subscription to ensure it meets
// requirements to have facts distributed to it.
func (subscription Subscription) validate() error {
	if subscription.Contact == "" {
		return errors.New("subscription requires contact")
	}

	return nil
}

func _CreateSubscription(writer interface{}, emitter *emission.Emitter) func(sub Subscription) {
	create := func(sub Subscription) (Subscription, error) {
		if err := sub.validate(); err != nil {
			return Subscription{}, err
		}
		return Subscription{}, nil
	}

	return emitResult[Subscription, Subscription, error](
		emitter,
		EventSubscriptionCreated,
		EventError,
	)(create)
}

func _DeleteSubscription(deleter interface{}, emitter *emission.Emitter) func(contact string) {
	delete := func(contact string) (interface{}, error) {
		return nil, nil
	}

	return emitResult[string, interface{}, error](
		emitter,
		EventSubscriptionDeleted,
		EventError,
	)(delete)
}

func _ListSubscriptions(lister interface{}, emitter *emission.Emitter) func(_ any) {
	list := func(_ any) ([]Subscription, error) {
		return nil, nil
	}

	return emitResult[any, []Subscription, error](
		emitter,
		EventSubscriptionsListed,
		EventError,
	)(list)
}
