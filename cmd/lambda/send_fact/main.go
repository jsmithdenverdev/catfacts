package main

import (
	"log"
	"os"

	"github.com/jsmithdenverdev/catfacts/catfactninja"
	"github.com/jsmithdenverdev/catfacts/dynamodb"
	"github.com/jsmithdenverdev/catfacts/lambda"
	"github.com/jsmithdenverdev/catfacts/twilio"
)

func main() {
	var (
		table       = os.Getenv("TABLE_NAME")
		from        = os.Getenv("TWILIO_FROM")
		sid         = os.Getenv("TWILIO_SID")
		token       = os.Getenv("TWILIO_TOKEN")
		storeLogger = log.New(os.Stdout, "[DYNAMO INFO] ", log.Ldate)
	)

	store, err := dynamodb.NewSubscriptionStore(table, storeLogger)

	if err != nil {
		panic(err)
	}

	smsSender := twilio.NewSMSSender(from, sid, token)
	retriever := catfactninja.CatFactNinjaClient{}

	handler := lambda.SendFactHandler{
		SubscriptionStore: store,
		FactRetriever:     retriever,
		SMSSender:         smsSender,
	}

	lambda.Start(handler.Handle)
}
