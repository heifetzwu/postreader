package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	// "github.com/aws/aws-lambda-go/lambda"
)

func main() {
	// lambda.Start(handler)
	// lambda.Start(getPostHandler)
	// lambda.Start(GetS3LambdaHandler)
	lambda.Start(getDB)
	// getS3()
}
