package lambda

import (
	"catfacts"
	"catfacts/twilio"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

func ManageSubscriptionHandler(service catfacts.SubscriptionService) func(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

func AuthorizerHandler(f catfacts.CredentialFetcher) func(request events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	return func(request events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
		// authorize the request
		authHeader := request.Headers["Authorization"]
		if authHeader == "" {
			return events.APIGatewayCustomAuthorizerResponse{}, catfacts.ErrMissingAuth
		}

		auth := strings.Fields(authHeader)
		if len(auth) != 2 {
			return events.APIGatewayCustomAuthorizerResponse{}, catfacts.ErrMalformedAuth
		}

		err := catfacts.AuthorizeBasicHeader(f, auth[1], "/twilio/auth")
		if err != nil {
			return events.APIGatewayCustomAuthorizerResponse{}, err
		}

		return events.APIGatewayCustomAuthorizerResponse{
			PrincipalID: "user",
			PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
				Version: "2012-10-17",
				Statement: []events.IAMPolicyStatement{
					{
						Action:   []string{"execute-api:Invoke"},
						Effect:   "Allow",
						Resource: []string{request.MethodArn},
					},
				},
			},
		}, nil
	}
}
