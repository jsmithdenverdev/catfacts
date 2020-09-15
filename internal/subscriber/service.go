package subscriber

import (
	"errors"
	"fmt"
)

type Service struct {
	store Store
}

//CreateSubscriber creates a new subscriber for a given contact. If a subscriber already
// exists with this contact an error will be returned.
func (s *Service) CreateSubscriber(contact Contact) error {
	if len(contact) == 0 {
		return errors.New("no contact supplied")
	}

	existing, err := s.store.Read(contact)

	if err != nil {
		return fmt.Errorf("could not check for existing subscriber: %w", err)
	}

	if len(existing.Contact) > 0 {
		return fmt.Errorf("a subscriber with this contact already exists: %s", contact)
	}

	err = s.store.Write(Subscriber{Contact: contact})

	if err != nil {
		return fmt.Errorf("could not write subscriber: %w", err)
	}

	return nil
}

//DeleteSubscriber deletes a subscriber for a given contact.
func (s *Service) DeleteSubscriber(contact Contact) error {
	err := s.store.Delete(contact)

	if err != nil {
		return fmt.Errorf("could not delete subscriber: %w", err)
	}

	return nil
}

//ListSubscribers returns a list of all subscribers.
func (s *Service) ListSubscribers() ([]Subscriber, error) {
	subscribers, err := s.store.List()

	if err != nil {
		return nil, fmt.Errorf("could not list subscribers: %w", err)
	}

	return subscribers, nil
}

func NewService(store Store) Service {
	return Service{
		store,
	}
}
