package facts

import (
	"errors"
	"fmt"
	"gitlab.com/jsmithdenverdev/catfacts/internal/distribution"
	"gitlab.com/jsmithdenverdev/catfacts/internal/subscriber"
)

type Fact = string

type Service struct {
	store       subscriber.Store
	distributor distribution.Distributor
	retriever   Retriever
}

func (s Service) DistributeToSubscribers() error {
	subscribers, err := s.store.List()

	if err != nil {
		return fmt.Errorf("could not retrieve list of subscribers: %w", err)
	}

	fact, err := s.retriever.Retrieve()

	if err != nil {
		return fmt.Errorf("could not retrieve fact: %w", err)
	}

	if len(fact) == 0 {
		return errors.New("retriever returned blank fact")
	}

	sendErrs := make([]error, 0)

	for _, sub := range subscribers {
		err := s.distributor.Distribute(sub.Contact, fmt.Sprintf("Meow! %v =^._.^=", fact))

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

func NewFactService(store subscriber.Store, distributor distribution.Distributor, retriever Retriever) Service {
	return Service{
		store,
		distributor,
		retriever,
	}
}
