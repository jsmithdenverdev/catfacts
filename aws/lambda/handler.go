package lambda

import (
	"catfacts"
	"catfacts/twilio"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func ManageSubscriptionHandler(service catfacts.SubscriptionService, credentialFetcher catfacts.CredentialFetcher) func(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		// authorize the request
		authHeader := event.Headers["Authorization"]
		if authHeader == "" {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusUnauthorized,
				Headers: map[string]string{
					"WWW-Authenticate": "Basic realm=\"Cat Facts\"",
				},
			}, catfacts.ErrUnauthorized
		}

		auth := strings.Fields(authHeader)
		if len(auth) != 2 {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusUnauthorized,
				Body:       catfacts.ErrMalformedAuth.Error(),
			}, catfacts.ErrMalformedAuth
		}

		err := catfacts.AuthorizeBasicHeader(credentialFetcher, auth[1], "/twilio/auth")
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusUnauthorized,
				Body:       catfacts.ErrMalformedAuth.Error(),
			}, err
		}

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

		reply, err := catfacts.ManageSMSSubscription(service, body, contact)

		if err != nil {
			return events.APIGatewayProxyResponse{
				Body:       err.Error(),
				StatusCode: http.StatusInternalServerError,
			}, err
		}

		x, err := twilio.GenerateTwiml(reply)
		if err != nil {
			return events.APIGatewayProxyResponse{
				Body:       err.Error(),
				StatusCode: http.StatusInternalServerError,
			}, err
		}

		return events.APIGatewayProxyResponse{
			Body:       string(x),
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"Content-Type": "application/xml",
			},
		}, nil
	}
}

func SendFactHandler(s catfacts.SubscriptionService, d catfacts.Distributor, fr catfacts.FactRetriever) func() error {
	return func() error {
		return catfacts.DistributeFactToSubscribers(s, d, fr)
	}
}
