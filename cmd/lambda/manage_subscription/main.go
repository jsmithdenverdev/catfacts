package main

import (
	"catfacts/aws/dynamo"
	handlers "catfacts/aws/lambda"
	"catfacts/subscription"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
)

func main() {
	table := os.Getenv("DYNAMODB_TABLE")
	store := dynamo.NewSubscriberStore(table)
	service := subscription.NewService(&store)

	handler := handlers.ManageSubscriptionHandler(service)
	lambda.Start(handler)
}
