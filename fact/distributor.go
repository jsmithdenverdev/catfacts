package fact

import (
	"catfacts/subscription"
	"errors"
	"strings"
)

type subscriberStore interface {
	All() ([]*subscription.Subscriber, error)
}

type sender interface {
	Send(to string, message string) error
}

type Distributor struct {
	store  subscriberStore
	sender sender
}

type Retriever = func() (string, error)

func NewDistributor(store subscriberStore, sender sender) Distributor {
	return Distributor{
		store,
		sender,
	}
}

func (d Distributor) DistributeFactToSubscribers(r Retriever) error {
	errs := make([]error, 0)
	fact, err := r()

	subscribers, err := d.store.All()
	if err != nil {
		return err
	}

	for _, sub := range subscribers {
		err := d.sender.Send(sub.Contact, fact)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		builder := strings.Builder{}
		for _, err := range errs {
			builder.WriteString(err.Error())
		}

		return errors.New(builder.String())
	}

	return nil
}
