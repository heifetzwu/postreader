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
	"github.com/aws/aws-sdk-go/service/sns"

	// "github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"
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

func newDB(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var err error
	// var expr expression.Expression
	var tablename, sndTopic string
	tableConst := "POSTS_TABLE"
	snsTopicConst := "SNS_TOPIC"

	body := Item{}
	// fmt.Println("before unmarshal", request.Body)
	err = json.Unmarshal([]byte(request.Body), &body)

	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, nil
	}
	uid := uuid.New()
	body.Id = uid.String()
	fmt.Println("body , uid= ", body, uid.String())

	// idstr := request.QueryStringParameters["postId"]

	// log.Println("ID = ", idstr)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)
	if os.Getenv(tableConst) == "" {
		tablename = "posts"
	} else {
		tablename = os.Getenv(tableConst)
	}
	fmt.Println("### tablename = ", tablename)

	if os.Getenv(snsTopicConst) == "" {
		sndTopic = "arn:aws:sns:ap-southeast-1:019907068212:postReader"
	} else {
		sndTopic = os.Getenv(snsTopicConst)
	}
	fmt.Println("## snsTopicConst = ", snsTopicConst)

	av, err := dynamodbattribute.MarshalMap(body)
	if err != nil {
		log.Fatalf("Got error marshalling new movie item: %s", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tablename),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}

	snsMsg := uid.String()

	svcSNS := sns.New(sess)
	resultSNS, err := svcSNS.Publish(&sns.PublishInput{
		Message:  &snsMsg,
		TopicArn: &sndTopic,
	})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("### messageid = ", *resultSNS.MessageId)
	res := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		// 	// Headers:    map[string]string{"Content-Type": "text/json; charset=utf-8"},
		Headers: map[string]string{
			"Content-Type":                 "application/text",
			"Access-Control-Allow-Headers": "Content-Type",
			"access-control-allow-origin":  "*",
			"Access-Control-Allow-Methods": "OPTIONS,POST,GET",
		},

		// 	// Headers: {
		// 	// 	"Content-Type":                "application/json",
		// 	// 	"access-control-allow-origin": "*"
		// 	// },
		Body: uid.String(),
	}
	return res, nil
}
