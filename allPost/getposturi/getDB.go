package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type Item struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	Text   string `json:"text"`
	Voice  string `json:"voice"`
	Url    string `json:"url"`
}

type InputType struct {
	Id string `json:"postId"`
}

// type Response struct {
// 	StatusCode int               `json:"statusCode"`
// 	Headers    map[string]string `json:"headers"`
// 	Body       string            `json:"body"`
// }

func getDB(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var err error
	var expr expression.Expression
	var tablename string
	tableConst := "POSTS_TABLE"

	idstr := request.QueryStringParameters["postId"]

	log.Println("ID = ", idstr)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)
	if os.Getenv(tableConst) == "" {
		tablename = "posts"
	} else {
		tablename = os.Getenv(tableConst)
	}
	fmt.Println("tablename = ", tablename)
	// idstr := "d93edd64-0fe7-4b51-b01d-8fafd7e294de"
	// idstr := in.Id
	proj := expression.NamesList(expression.Name("id"), expression.Name("status"), expression.Name("text"), expression.Name("url"), expression.Name("voice"))
	// var expr interface{}

	if idstr != "*" {
		filt := expression.Name("id").Equal(expression.Value(idstr))
		expr, err = expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	} else {
		expr, err = expression.NewBuilder().WithProjection(proj).Build()
	}

	// expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	if err != nil {
		log.Fatalf("Got error building expression: %s", err)
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tablename),
	}

	result, err := svc.Scan(params)

	if err != nil {
		log.Fatalf("Query API call failed: %s", err)
	}

	numItems := 0

	// fmt.Println("result.Items = %s", result.Items)

	allitems := []Item{}

	// err = dynamodbattribute.UnmarshalMap(result.Items, allitems)

	// fmt.Println("allItems = %s", allitems)

	for _, i := range result.Items {

		item := Item{}

		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			log.Fatalf("Got error unmarshalling: %s", err)
		}

		numItems++

		fmt.Println("Id: ", item.Id)
		fmt.Println("Status:", item.Status)
		fmt.Println("Text:", item.Text)
		fmt.Println("Voice:", item.Voice)
		fmt.Println("URL:", item.Url)
		fmt.Println()
		allitems = append(allitems, item)
	}

	fmt.Println("Found items:", numItems)
	// fmt.Println("all items:", allitems)
	// return Response{
	// 		StatusCode: 200,
	// 		Headers:    map[string]string{"Content-Type": "application/json"},
	// 		Body:       fmt.Sprintf("%s", (allitems)),
	// 	},
	// 	nil

	alljson, err := json.Marshal(allitems)
	if err != nil {
		log.Fatalf("marshal error %s", err)

	}

	// return fmt.Sprintf("%s", (allitems)), nil

	res := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		// Headers:    map[string]string{"Content-Type": "text/json; charset=utf-8"},
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"access-control-allow-origin": "*",
		},

		// Headers: {
		// 	"Content-Type":                "application/json",
		// 	"access-control-allow-origin": "*"
		// },
		Body: string(alljson),
	}
	return res, nil
}
