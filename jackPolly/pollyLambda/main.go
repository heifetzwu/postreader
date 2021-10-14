package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	// "github.com/aws/aws-lambda-go/lambda"
)

var (
	// DefaultHTTPGetAddress Default Address
	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"

	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("No IP in HTTP response")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("Non 200 Response found")
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp, err := http.Get(DefaultHTTPGetAddress)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if resp.StatusCode != 200 {
		return events.APIGatewayProxyResponse{}, ErrNon200Response
	}

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if len(ip) == 0 {
		return events.APIGatewayProxyResponse{}, ErrNoIP
	}

	fmt.Println("ip = ", string(ip))

	m := make(map[string]string)
	m["one"] = "2-01"
	m["two"] = "2-02"

	jsonBytes, err := json.Marshal(m)

	rstr := fmt.Sprintf(string(jsonBytes))
	fmt.Println("jsonBytes", jsonBytes)

	return events.APIGatewayProxyResponse{
		Body: rstr,
		// Body:       "test" + string(jsonBytes),
		// Body:       "test",
		StatusCode: 200,
	}, nil
}
func main() {
	// lambda.Start(handler)
	// lambda.Start(getPostHandler)
	// lambda.Start(GetS3LambdaHandler)
	lambda.Start(getDB)
	// getS3()
}
