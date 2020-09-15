package app

import (
	"fmt"
	"gitlab.com/jsmithdenverdev/catfacts/internal/subscriber"
	"os"
)

type App struct {
	subscriberService subscriber.SubscriberService
}

func CreateApp() (App, error) {
	dataSource := os.Getenv("DATA_SOURCE")

	// create a new sqlite subscriber store
	store, err := subscriber.NewSqliteSubscriberStore(dataSource)

	if err != nil {
		return App{}, fmt.Errorf("could not create App: %w", err)
	}

	// create services
	subscriberService := subscriber.NewSubscriberService(store)

	// create App
	return App{subscriberService}, nil
}
