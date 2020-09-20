package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jsmithdenverdev/catfacts/internal/aws/dynamodb"
	"github.com/jsmithdenverdev/catfacts/internal/fact"
	"github.com/jsmithdenverdev/catfacts/internal/twilio"
	"log"
	"os"
)

func handler() error {
	twilioSid := os.Getenv("TWILIO_SID")
	twilioToken := os.Getenv("TWILIO_TOKEN")
	twilioFrom := os.Getenv("TWILIO_FROM")
	table := os.Getenv("DYNAMODB_TABLE")

	store := dynamodb.NewSubscriberStore(table)
	sender := twilio.NewFactSender(twilioSid, twilioToken, twilioFrom)
	service := fact.NewService(store, sender)

	err := service.SendFactToSubscribers(fact.RetrieveFromApi)

	if err != nil {
		log.Print(err.Error())
	}

	return err
}

func main() {
	lambda.Start(handler)
}
