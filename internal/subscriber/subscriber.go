package subscriber

// Contact is a string representing a phone number.
type Contact = string

// Subscriber represents a subscriber for subscriptions. A subscriber has a single field. Contact which represents the
// subscribers phone number.
type Subscriber struct {
	Contact Contact `dynamodbav:"contact"`
}

// Store is an interface that provides crud operations for Subscriber.
type Store interface {
	Read(contact Contact) (Subscriber, error)
	Write(subscriber Subscriber) error
	List() ([]Subscriber, error)
	Delete(contact Contact) error
}
