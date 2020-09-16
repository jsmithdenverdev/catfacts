package main

import (
	"fmt"
	"github.com/jsmithdenverdev/catfacts/internal/postgres"
	"github.com/jsmithdenverdev/catfacts/internal/subscriber"
	"os"
)

type app struct {
	subscriberService subscriber.Service
}

func createApp() (app, error) {
	databaseUrl := os.Getenv("DATABASE_URL")

	// create a new subscriber store
	store, err := postgres.NewSubscriberStore(databaseUrl)

	if err != nil {
		return app{}, fmt.Errorf("could not connect to subscriber store: %w", err)
	}

	// create services
	subscriberService := subscriber.NewService(store)

	// create App
	return app{subscriberService}, nil
}
