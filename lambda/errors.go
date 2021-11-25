package lambda

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func APIGatewayProxyResponseError(err error) (events.APIGatewayProxyResponse, error) {
	bytes, encodingErr := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: err.Error(),
	})

	if encodingErr != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
			Headers: map[string]string{
				"Content-Type": "text/plain",
			},
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       string(bytes),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, err
}
