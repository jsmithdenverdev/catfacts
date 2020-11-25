package main

import (
	"catfacts/aws/dynamo"
	handlers "catfacts/aws/lambda"
	"catfacts/aws/ssm"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"os"
)

func main() {
	table := os.Getenv("DYNAMODB_TABLE")

	sess := session.Must(session.NewSession())
	service := dynamo.NewSubscriptionService(sess, table)
	credentialFetcher := ssm.NewSsmCredentialFetcher(sess)

	handler := handlers.ManageSubscriptionHandler(&service, credentialFetcher)
	lambda.Start(handler)
}
