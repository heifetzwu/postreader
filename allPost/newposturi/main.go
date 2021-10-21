package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	// "github.com/aws/aws-lambda-go/lambda"
)

func main() {
	// lambda.Start(handler)
	// lambda.Start(getPostHandler)
	// lambda.Start(GetS3LambdaHandler)
	fmt.Println("newPost start")
	lambda.Start(newDB)
	// getS3()
}
