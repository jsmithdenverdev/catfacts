package main

import (
	"catfacts/aws/dynamo"
	handlers "catfacts/aws/lambda"
	"catfacts/fact"
	"catfacts/twilio"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
)

func main() {
	sid := os.Getenv("TWILIO_SID")
	token := os.Getenv("TWILIO_TOKEN")
	from := os.Getenv("TWILIO_FROM")
	table := os.Getenv("DYNAMODB_TABLE")

	store := dynamo.NewSubscriberStore(table)
	sender := twilio.NewSmsSender(sid, token, from)
	distributor := fact.NewDistributor(&store, sender)

	handler := handlers.SendFactHandler(distributor)
	lambda.Start(handler)
}
