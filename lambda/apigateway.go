package lambda

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

type APIGatewayHandler interface {
	Handle(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}
