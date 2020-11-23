package subscription

type store interface {
	Insert(sub Subscriber) error
	Delete(contact string) error
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

func (s Service) ManageSMSSubscription(body, contact string) (Reply, error) {
	action := GetAction(body)

	switch action {
	case ActionSubscribe:
		subscriber := Subscriber{
			Contact: contact,
		}
		if err := s.store.Insert(subscriber); err != nil {
			if err == ErrSubscriptionExists {
				return ReplyAlreadyExists, nil
			} else {
				return "", err
			}
		} else {
			return ReplyWelcome, nil
		}
	case ActionUnsubscribe:
		if err := s.store.Delete(contact); err != nil {
			return "", nil
		} else {
			return ReplyGoodbye, nil
		}
	case ActionHelp:
		return ReplyHelp, nil
	case ActionUnknown:
		return ReplyUnknown, nil
	}

	return ReplyUnknown, nil
}
