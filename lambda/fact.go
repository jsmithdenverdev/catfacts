package lambda

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jsmithdenverdev/catfacts"
)

type SendFactHandler struct {
	SubscriptionStore catfacts.SubscriptionStore
	FactRetriever     catfacts.FactRetriever
	SMSSender         catfacts.SMSSender
}

func (handler SendFactHandler) Handle(ctx context.Context, event events.CloudWatchEvent) error {
	fact, err := handler.FactRetriever.Retrieve()

	if err != nil {
		return err
	}

	subscriptions, err := handler.SubscriptionStore.All(ctx)

	if err != nil {
		return err
	}

	for _, subscription := range subscriptions {
		if err := handler.SMSSender.Send(catfacts.SMS{
			Body: fmt.Sprintf("%s %s", catfacts.SMSBodyFact, fact),
			To:   subscription.Contact,
		}); err != nil {
			return err
		}
	}

	return nil
}
