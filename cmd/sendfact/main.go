package main

import (
	"fmt"
	"github.com/jsmithdenverdev/catfacts/internal/fact"
	"github.com/jsmithdenverdev/catfacts/internal/sqlite3"
	"github.com/jsmithdenverdev/catfacts/internal/twilio"
	"log"
	"os"
)

type app struct {
	factService fact.Service
}

func createApp() (app, error) {
	twilioSid := os.Getenv("TWILIO_SID")
	twilioToken := os.Getenv("TWILIO_TOKEN")
	twilioFrom := os.Getenv("TWILIO_FROM")
	dataSource := os.Getenv("DATA_SOURCE")

	// create a new sqlite subscriber store
	store, err := sqlite3.NewSubscriberStore(dataSource)

	if err != nil {
		return app{}, fmt.Errorf("could not create app: %w", err)
	}

	sender := twilio.NewFactSender(twilioSid, twilioToken, twilioFrom)
	factService := fact.NewService(store, sender)

	// create app
	return app{factService}, nil
}

func run() error {
	app, err := createApp()

	if err != nil {
		return fmt.Errorf("could not create app: %w", err)
	}

	err = app.factService.SendFactToSubscribers(fact.RetrieveFactFromApi)

	if err != nil {
		return fmt.Errorf("could not send fact to subscribers: %w", err)
	}

	return nil
}

func main() {
	err := run()

	if err != nil {
		log.Fatal(err)
	}
}
