package lambda

import (
	"context"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jsmithdenverdev/catfacts"
	"github.com/jsmithdenverdev/catfacts/twilio/twiml"
)

type TwilioCallbackHandler struct {
	SubscriptionStore catfacts.SubscriptionStore
}

func (handler TwilioCallbackHandler) Handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	query, err := url.ParseQuery(req.Body)

	if err != nil {
		return APIGatewayProxyResponseError(err)
	}

	from := query.Get("From")
	body := query.Get("Body")

	sms := catfacts.SMS{
		From: from,
		Body: body,
	}

	command := sms.ToCommand()

	switch command {
	case catfacts.CommandSubscribe:
		return handler.handleSubscribe(ctx, sms)
	case catfacts.CommandUnsubscribe:
		return handler.handleUnsubscribe(ctx, sms)
	case catfacts.CommandCuss:
		return handler.handleCuss(sms)
	case catfacts.CommandHelp:
		fallthrough
	default:
		return handler.handleHelp(sms)
	}
}

func (handler TwilioCallbackHandler) handleSubscribe(ctx context.Context, sms catfacts.SMS) (events.APIGatewayProxyResponse, error) {
	subscription := catfacts.Subscription{
		Contact: sms.From,
	}

	if err := handler.SubscriptionStore.Insert(ctx, subscription); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return twimlResponse(catfacts.SMSBodyWelcome)
}

func (handler TwilioCallbackHandler) handleUnsubscribe(ctx context.Context, sms catfacts.SMS) (events.APIGatewayProxyResponse, error) {
	if err := handler.SubscriptionStore.Delete(ctx, sms.From); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return twimlResponse(catfacts.SMSBodyGoodbye)
}

func (handler TwilioCallbackHandler) handleHelp(sms catfacts.SMS) (events.APIGatewayProxyResponse, error) {
	return twimlResponse(catfacts.SMSBodyHelp)
}

func (handler TwilioCallbackHandler) handleCuss(sms catfacts.SMS) (events.APIGatewayProxyResponse, error) {
	return twimlResponse(catfacts.SMSBodyCuss)
}

// twimlResponse marshalls the provided body into Twiml and attaches it to an
// APIGatewayProxyResponse as the Body. The response will have a Content-Type
// header of text/xml. If an error occurs during marshalling, an APIGatewayProxyResponse
// containing the marhsalling error in the Body and a Content-Type header of text/plain
// will be returned. If an error occured during marshalling the error is returned
// as the second item in the response tuple.
func twimlResponse(body string) (events.APIGatewayProxyResponse, error) {
	// Marshal Twiml response body
	responseBytes, encodingErr := twiml.Marshal(body)

	// If marshalling the response failed send a 500 with the failure.
	// Send the original error as the second value to enable reporting on it.
	if encodingErr != nil {
		return events.APIGatewayProxyResponse{}, encodingErr
	}

	// Send Twilio a 200 response with the Twiml body. This will be sent back to
	// the original SMS sender.
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBytes),
		Headers: map[string]string{
			"Content-Type": "text/xml",
		},
	}, nil
}
