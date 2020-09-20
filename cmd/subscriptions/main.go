package main

import (
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jsmithdenverdev/catfacts/internal/aws/dynamodb"
	"github.com/jsmithdenverdev/catfacts/internal/subscriber"
	"github.com/jsmithdenverdev/catfacts/internal/twilio"
	"log"
	"net/http"
	"net/url"
	"os"
)

func handler(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var err error
	table := os.Getenv("DYNAMODB_TABLE")
	store := dynamodb.NewSubscriberStore(table)
	service := subscriber.NewService(store)

	params, err := url.ParseQuery(event.Body)
	if err != nil {
		log.Print(err.Error())

		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	body := params.Get("Body")
	contact := params.Get("From")

	if len(contact) == 0 {
		err = errors.New("could not create subscriber: malformed twilio request: missing parameter From")

		log.Print(err.Error())

		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: http.StatusBadRequest,
		}, err
	}

	switch twilio.GetTwilioAction(body) {
	case twilio.OptIn:
		// create the subscriber
		err = service.CreateSubscriber(contact)

		// return an error if creation failed
		if err != nil {
			log.Print(err.Error())
			return events.APIGatewayProxyResponse{
				Body:       err.Error(),
				StatusCode: http.StatusInternalServerError,
			}, err
		}

		resp, err := twilio.GenerateTwiml("Meow! Welcome to CatFacts! =^._.^=. Cancel your subscription at any time by replying STOP.")

		if err != nil {
			log.Print(err.Error())
			return events.APIGatewayProxyResponse{
				Body:       err.Error(),
				StatusCode: http.StatusInternalServerError,
			}, err
		}

		return events.APIGatewayProxyResponse{
			Body:       string(resp),
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"Content-Type": "text/xml",
			},
		}, nil
	case twilio.OptOut:
		// delete the subscriber
		err = service.DeleteSubscriber(contact)

		// return an error if deleting failed
		if err != nil {
			log.Print(err.Error())
			return events.APIGatewayProxyResponse{
				Body:       err.Error(),
				StatusCode: http.StatusInternalServerError,
			}, err
		}

		// write a twiml response
		resp, err := twilio.GenerateTwiml("Meow! Farewell from CatFacts! =^._.^=")

		if err != nil {
			log.Print(err.Error())
			return events.APIGatewayProxyResponse{
				Body:       err.Error(),
				StatusCode: http.StatusInternalServerError,
				Headers: map[string]string{
					"Content-Type": "text/xml",
				},
			}, err
		}

		return events.APIGatewayProxyResponse{
			Body:       string(resp),
			StatusCode: http.StatusOK,
		}, nil
	}

	// create and return an invalid operation error
	resp, err := twilio.GenerateTwiml("Meow! That request was not understood. =^._.^=")

	if err != nil {
		log.Print(err.Error())
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(resp),
		StatusCode: http.StatusBadRequest,
		Headers: map[string]string{
			"Content-Type": "text/xml",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
