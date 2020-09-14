package subscriber

import (
	"fmt"
)

type Contact = string

type Subscriber struct {
	Contact Contact
}

type Service struct {
	store Store
}

type Store interface {
	Read(contact string) (*Subscriber, error)
	Write(subscriber Subscriber) error
	List() ([]Subscriber, error)
	Delete(contact string) error
}

func NewSubscriberService(store Store) Service {
	return Service{
		store,
	}
}

func (s Service) CreateSubscriber(contact Contact) error {
	existing, err := s.store.Read(contact)

	if err != nil {
		return fmt.Errorf("could not check for existing subscriber: %w", err)
	}

	if existing != nil {
		return fmt.Errorf("a subscriber with this contact already exists: %s", contact)
	}

	subscriber := Subscriber{Contact: contact}
	err = s.store.Write(subscriber)

	if err != nil {
		return fmt.Errorf("could not Write subscriber: %w", err)
	}

	return nil
}

func (s Service) DeleteSubscriber(contact string) error {
	err := s.store.Delete(contact)

	if err != nil {
		return fmt.Errorf("could not Delete subscriber: %w", err)
	}

	return nil
}
