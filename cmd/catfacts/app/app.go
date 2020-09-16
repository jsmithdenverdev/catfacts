package app

import (
	"fmt"
	"github.com/jsmithdenverdev/catfacts/internal/sqlite3"
	"github.com/jsmithdenverdev/catfacts/internal/subscriber"
	"os"
)

type App struct {
	subscriberService subscriber.Service
}

func CreateApp() (App, error) {
	dataSource := os.Getenv("DATA_SOURCE")

	// create a new sqlite subscriber store
	store, err := sqlite3.NewSubscriberStore(dataSource)

	if err != nil {
		return App{}, fmt.Errorf("could not create App: %w", err)
	}

	// create services
	subscriberService := subscriber.NewService(store)

	// create App
	return App{subscriberService}, nil
}
