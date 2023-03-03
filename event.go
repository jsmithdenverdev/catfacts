package catfacts

type Event = string

const (
	EventError                 Event = "ERROR"
	EventSubscriptionCreated   Event = "SUBSCRIPTION_CREATED"
	EventSubscriptionDeleted   Event = "SUBSCRIPTION_DELETE"
	EventSubscriptionsListed   Event = "SUBSCRIPTIONS_LISTED"
	EventDistributionSucceeded Event = "DISTRIBUTION_SUCCEEDED"
	EventCussedOut             Event = "CUSSED_OUT"
)
