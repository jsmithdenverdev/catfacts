package lambda

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

type APIGatewayLoggingHandler struct {
	Logger *log.Logger
	Inner  APIGatewayHandler
}

func (handler APIGatewayLoggingHandler) Handle(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// print the raw request JSON
	json.NewEncoder(log.Writer()).Encode(event)

	return handler.Inner.Handle(ctx, event)
}
