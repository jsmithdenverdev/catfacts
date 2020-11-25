package main

import (
	"catfacts/aws/dynamo"
	handlers "catfacts/aws/lambda"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
)

func main() {
	table := os.Getenv("DYNAMODB_TABLE")

	sess := session.Must(session.NewSession())
	service := dynamo.NewSubscriptionService(sess, table)

	handler := handlers.ManageSubscriptionHandler(&service)
	lambda.Start(handler)
}
