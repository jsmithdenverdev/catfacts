package fact

import "github.com/jsmithdenverdev/catfacts/internal/subscriber"

type Fact = string

type Retriever func() (Fact, error)

type Sender interface {
	Send(to subscriber.Contact, fact Fact) error
}
