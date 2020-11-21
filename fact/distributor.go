package fact

import (
	"catfacts/subscription"
	"errors"
	"strings"
)

type subscriberLister interface {
	ListAll() ([]*subscription.Subscriber, error)
}

type sender interface {
	Send(to string, message string) error
}

type Distributor struct {
	lister subscriberLister
	sender sender
}

type FactRetriever = func() (string, error)

func NewDistributor(lister subscriberLister, sender sender) Distributor {
	return Distributor{
		lister,
		sender,
	}
}

func (d Distributor) DistributeFactToSubscribers(r FactRetriever) error {
	errs := make([]error, 0)
	fact, err := r()

	subscribers, err := d.lister.ListAll()
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
