package main

import (
	handlers "catfacts/aws/lambda"
	"catfacts/aws/ssm"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
)

func main() {

	sess := session.Must(session.NewSession())
	fetcher := ssm.NewSsmCredentialFetcher(sess)

	handler := handlers.AuthorizerHandler(&fetcher)
	lambda.Start(handler)
}
