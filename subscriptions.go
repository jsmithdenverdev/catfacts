package catfacts

import (
	"context"
	"errors"
)

type Subscription struct {
	Contact string
}

type SubscriptionStore interface {
	Insert(ctx context.Context, subscription Subscription) error
	Delete(ctx context.Context, contact string) error
	All(ctx context.Context) ([]Subscription, error)
}

// Validate performs basic validation on a Subscription to ensure it meets
// requirements to have facts distributed to it.
func (subscription Subscription) Validate() error {
	if subscription.Contact == "" {
		return errors.New("subscription requires contact")
	}

	return nil
}
