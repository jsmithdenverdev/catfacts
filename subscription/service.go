package subscription

import (
	"errors"
	"fmt"
)

type store interface {
	Insert(sub Subscriber) error
	Delete(contact string) error
	All() ([]*Subscriber, error)
}

type Service struct {
	store store
}

type Error struct {
	message string
}

func (e Error) Error() string {
	return e.message
}

var (
	ErrSubscriptionExists = Error{
		message: "a subscription already exists for this contact",
	}
)

func NewService(store store) Service {
	return Service{
		store,
	}
}

func (s Service) Create(contact string) error {
	if contact == "" {
		return errors.New("empty contact supplied")
	}

	subscriber := Subscriber{
		Contact: contact,
	}

	return s.store.Insert(subscriber)
}

func (s Service) Delete(contact string) error {
	if contact == "" {
		return errors.New("empty contact supplied")
	}

	err := s.store.Delete(contact)

	if err != nil {
		return fmt.Errorf("could not delete subscription from store: %w", err)
	}

	return nil
}

func (s Service) ListAll() ([]*Subscriber, error) {
	results, err := s.store.All()

	if err != nil {
		return nil, fmt.Errorf("could not list subscriptions from store: %w", err)
	}

	return results, nil
}
