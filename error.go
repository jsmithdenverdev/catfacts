package catfacts

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrSubscriptionExists = Error("a subscription already exists for this contact")
	ErrUnauthorized       = Error("unauthorized")
	ErrMissingAuth        = Error("missing authorization header")
	ErrMalformedAuth      = Error("malformed authorization header")
)
