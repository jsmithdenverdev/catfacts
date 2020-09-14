package app

import (
	"fmt"
	"gitlab.com/jsmithdenverdev/catfacts/internal/distribution"
	"gitlab.com/jsmithdenverdev/catfacts/internal/facts"
	"gitlab.com/jsmithdenverdev/catfacts/internal/subscriber"
	"os"
)

type App struct {
	SubscriberService subscriber.Service
	FactService       facts.Service
}

func CreateApp() (App, error) {
	twilioSid := os.Getenv("TWILIO_SID")
	twilioToken := os.Getenv("TWILIO_TOKEN")
	twilioFrom := os.Getenv("TWILIO_FROM")

	// create a new sqlite subscriber store
	store, err := subscriber.NewSqliteSubscriberStore("file:catfacts.db")

	if err != nil {
		return App{}, fmt.Errorf("could not create app: %w", err)
	}

	// create a retriever and distributor
	factRetriever := facts.NewCatFactNinjaRetriever()
	distributor := distribution.NewTwilioDistributor(twilioSid, twilioToken, twilioFrom)

	// create services
	subscriberService := subscriber.NewSubscriberService(store)
	factService := facts.NewFactService(store, distributor, factRetriever)

	// create app
	return App{
		SubscriberService: subscriberService,
		FactService:       factService,
	}, nil
}
