package main

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

// type Item struct {
// 	id string,
// }

// var tablename string  "aws-python-postreader-dev"

func getPostHandlerTBD(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	var buf bytes.Buffer
	body, err := json.Marshal(map[string]interface{}{
		"message": "Go serverless v1.0",
	})

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Context-Type":           "application/json",
			"X-MyCompany-Func-Reply": "hello-handler",
		},
	}
	return resp, nil
}
