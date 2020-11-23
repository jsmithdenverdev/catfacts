package lambda

import (
	"catfacts/fact"
	"catfacts/http/rest"
	"catfacts/subscription"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"net/http"
	"net/url"
)

func ManageSubscriptionHandler(service subscription.Service) func(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

		reply, err := service.ManageSMSSubscription(body, contact)

		if err != nil {
			return events.APIGatewayProxyResponse{
				Body:       err.Error(),
				StatusCode: http.StatusInternalServerError,
			}, err
		}

		x, err := rest.GenerateTwiml(reply)
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

func SendFactHandler(d fact.Distributor) func() error {
	return func() error {
		return d.DistributeFactToSubscribers(fact.RetrieveFactFromCatfactNinja)
	}
}
