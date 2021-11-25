package main

import (
	"log"
	"os"

	"github.com/jsmithdenverdev/catfacts/dynamodb"
	"github.com/jsmithdenverdev/catfacts/lambda"
)

func main() {
	var (
		table        = os.Getenv("TABLE_NAME")
		storeLogger  = log.New(os.Stdout, "[DYNAMO INFO] ", log.Ldate)
		lambdaLogger = log.New(os.Stdout, "[LAMBDA] [INFO] ", log.Ldate)
	)

	store, err := dynamodb.NewSubscriptionStore(table, storeLogger)

	if err != nil {
		panic(err)
	}

	var handler lambda.APIGatewayHandler

	handler = lambda.TwilioCallbackHandler{
		SubscriptionStore: store,
	}

	// Wrap the handler with trivial logging
	handler = lambda.APIGatewayLoggingHandler{
		Logger: lambdaLogger,
		Inner:  handler,
	}

	lambda.Start(handler.Handle)
}
