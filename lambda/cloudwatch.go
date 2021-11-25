package lambda

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

type CloudWatchEventHandler interface {
	Handle(ctx context.Context, event events.CloudWatchEvent) error
}
