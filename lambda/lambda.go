package lambda

import "github.com/aws/aws-lambda-go/lambda"

func Start(handler interface{}) {
	lambda.Start(handler)
}
