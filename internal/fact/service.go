package fact

import (
	"errors"
	"fmt"
	"github.com/jsmithdenverdev/catfacts/internal/subscriber"
)

type Service struct {
	store  subscriber.Store
	sender Sender
}

func (s *Service) SendFactToSubscribers(retriever Retriever) error {
	subs, err := s.store.List()

	if err != nil {
		return fmt.Errorf("could not get list of subscribers: %w", err)
	}

	fact, err := retriever()

	if err != nil {
		return fmt.Errorf("could not retrieve fact: %w", err)
	}

	sendErrs := make([]error, 0)

	for _, sub := range subs {
		err := s.sender.Send(sub.Contact, fact)

		if err != nil {
			sendErrs = append(sendErrs, err)
		}
	}

	if len(sendErrs) > 0 {
		composite := ""

		for _, err := range sendErrs {
			composite += err.Error()
			composite += "\n"
		}

		return errors.New(composite)
	}

	return nil
}

func NewService(store subscriber.Store, sender Sender) Service {
	return Service{
		store,
		sender,
	}
}
