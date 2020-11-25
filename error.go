package catfacts

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrSubscriptionExists = Error("a subscription already exists for this contact")
	ErrUnauthorized       = Error("unauthorized")
	ErrMalformedAuth      = Error("malformed authorization header")
)
