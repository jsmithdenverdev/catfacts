package app

import (
	"fmt"
	"gitlab.com/jsmithdenverdev/catfacts/internal"
	"os"
)

type App struct {
	subscriberService internal.SubscriberService
}

func CreateApp() (App, error) {
	dataSource := os.Getenv("DATA_SOURCE")

	// create a new sqlite subscriber store
	store, err := internal.NewSqliteSubscriberStore(dataSource)

	if err != nil {
		return App{}, fmt.Errorf("could not create App: %w", err)
	}

	// create services
	subscriberService := internal.NewSubscriberService(store)

	// create App
	return App{subscriberService}, nil
}
