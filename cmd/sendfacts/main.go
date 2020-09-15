package main

import (
	"fmt"
	"gitlab.com/jsmithdenverdev/catfacts/internal"
	"gitlab.com/jsmithdenverdev/catfacts/internal/subscriber"
	"log"
	"os"
)

type app struct {
	factService internal.FactService
}

func createApp() (app, error) {
	twilioSid := os.Getenv("TWILIO_SID")
	twilioToken := os.Getenv("TWILIO_TOKEN")
	twilioFrom := os.Getenv("TWILIO_FROM")
	dataSource := os.Getenv("DATA_SOURCE")

	// create a new sqlite subscriber store
	store, err := subscriber.NewSqliteSubscriberStore(dataSource)

	if err != nil {
		return app{}, fmt.Errorf("could not create app: %w", err)
	}

	loader := internal.NewApiFactLoader()
	distributor := internal.NewTwilioDistributor(twilioSid, twilioToken, twilioFrom)

	// create services
	factService := internal.NewFactService(store, loader, distributor)

	// create app
	return app{factService}, nil
}

func run() error {
	app, err := createApp()

	if err != nil {
		return fmt.Errorf("could not create app: %w", err)
	}

	err = app.factService.SendFactToSubscribers()

	if err != nil {
		return fmt.Errorf("could not send fact to subscribers")
	}

	return nil
}

func main() {
	err := run()

	if err != nil {
		log.Fatal(err)
	}
}
