package main

import (
	"catfacts"
	"catfacts/aws/dynamo"
	handlers "catfacts/aws/lambda"
	"catfacts/twilio"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"os"
)

func main() {
	sid := os.Getenv("TWILIO_SID")
	token := os.Getenv("TWILIO_TOKEN")
	from := os.Getenv("TWILIO_FROM")
	table := os.Getenv("DYNAMODB_TABLE")

	sess := session.Must(session.NewSession())
	service := dynamo.NewSubscriptionService(sess, table)
	distributor := twilio.NewSmsDistributor(sid, token, from)

	handler := handlers.SendFactHandler(&service, distributor, catfacts.RetrieveFactFromCatfactNinja)
	lambda.Start(handler)
}
