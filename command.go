package catfacts

import (
	"github.com/chuckpreslar/emission"
)

type Command string

const (
	CommandSubscribe   Command = "subscribe"
	CommandUnsubscribe Command = "unsubscribe"
	CommandCuss        Command = "cuss"
	CommandDistribute  Command = "distribute"
)

func HandleSubscribe(writer interface{}, emitter *emission.Emitter) {
	subscribe := func(sub Subscription) (Subscription, error) {
		if err := sub.validate(); err != nil {
			return Subscription{}, err
		}
		return Subscription{}, nil
	}

	handler := emitResult[Subscription, Subscription, error](
		emitter,
		EventSubscriptionCreated,
		EventError,
	)(subscribe)

	emitter.On(CommandSubscribe, handler)
}

func HandleUnsubscribe(writer interface{}, emitter *emission.Emitter) {
	unsubscribe := func(sub Subscription) (Subscription, error) {
		if err := sub.validate(); err != nil {
			return Subscription{}, err
		}
		return Subscription{}, nil
	}

	handler := emitResult[Subscription, Subscription, error](
		emitter,
		EventSubscriptionDeleted,
		EventError,
	)(unsubscribe)

	emitter.On(CommandUnsubscribe, handler)
}

func HandleCuss(sender interface{}, emitter *emission.Emitter) {
	cuss := func(_ any) (any, error) {
		return nil, nil
	}

	handler := emitResult[any, any, error](
		emitter,
		EventCussedOut,
		EventError,
	)(cuss)

	emitter.On(CommandCuss, handler)
}

func HandleDistribute(lister interface{}, retriever interface{}, sender interface{}, emitter *emission.Emitter) {
	distribute := func(_ any) (any, error) {
		return nil, nil
	}

	handler := emitResult[any, any, error](
		emitter,
		EventDistributionSucceeded,
		EventError,
	)(distribute)

	emitter.On(CommandDistribute, handler)
}
