package catfacts

type Event string

const (
	EventError          Event = "ERROR"
	SubscriptionCreated Event = "SUBSCRIPTION_CREATED"
	SubscriptionDeleted Event = "SUBSCRIPTION_DELETE"
)
