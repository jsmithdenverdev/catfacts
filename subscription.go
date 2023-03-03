package catfacts

import (
	"errors"
	"github.com/chuckpreslar/emission"
)

type Subscription struct {
	Contact string
}

// Validate performs basic validation on a Subscription to ensure it meets
// requirements to have facts distributed to it.
func (subscription Subscription) Validate() error {
	if subscription.Contact == "" {
		return errors.New("subscription requires contact")
	}

	return nil
}

func Result[TIn, TOk, TErr any](delegate func(TIn) (TOk, TErr), ok func(TOk), err func(TErr)) func(TIn) {
	return func(in TIn) {
		result, e := delegate(in)
		if e != nil {
			err(e)
			return
		}
		ok(result)
	}
}

func WithErrorEmitter[TIn, TOk, TErr any](emitter *emission.Emitter) func(delegate func(TIn) (TOk, TErr), ok func(TOk)) func(TIn) {
	return func(delegate func(TIn) (TOk, TErr), ok func(TOk)) func(TIn) {
		return Result[TIn, TOk, TErr](delegate, ok, func(err TErr) {
			emitter.Emit(EventError, err)
		})
	}
}

func LambdaSample() {
	e := emission.NewEmitter()
	w := new(bool)
	e.On(CommandSubscribe, WithErrorEmitter[Subscription, Subscription, error](e)(CreateSubscription(w), func(subscription Subscription) {
		e.Emit(SubscriptionCreated, subscription)
	}))
}

func withBehavior[TIn any, TSuccess any](
	onSuccess func(output TSuccess),
	onFailure func(err error),
	callback func(in TIn) (TSuccess, error)) func(in TIn) {
	return func(in TIn) {
		result, err := callback(in)
		if err != nil {
			onFailure(err)
		} else {
			onSuccess(result)
		}
	}
}

func CreateSubscription(writer interface{}) func(sub Subscription) (Subscription, error) {
	return func(sub Subscription) (Subscription, error) {
		return Subscription{}, nil
	}
}

func createSubscription(writer interface{}, emitter *emission.Emitter) func(subscription Subscription) {

	onSuccess := func(sub Subscription) {
		emitter.Emit(SubscriptionCreated, sub)
	}

	onFail := func(e error) {
		emitter.Emit(EventError, e)
	}

	create := func(sub Subscription) (Subscription, error) {
		return Subscription{}, nil
	}

	return withBehavior[Subscription, Subscription](onSuccess, onFail, create)
}

func DeleteSubscription(deleter interface{}) func(contact string) error {
	return func(contact string) error {
		return nil
	}
}

func ListSubscriptions(lister interface{}) func() ([]Subscription, error) {
	return func() ([]Subscription, error) {
		return nil, nil
	}
}
