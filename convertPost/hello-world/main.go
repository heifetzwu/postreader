package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Item struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	Text   string `json:"text"`
	Voice  string `json:"voice"`
	Url    string `json:"url"`
}

func handler(ctx context.Context, snsEvent events.SNSEvent) {
	var tablename string
	var idstr string
	tableConst := "POSTS_TABLE"
	// for _, record := range snsEvent.Records {
	// 	snsRecord := record.SNS
	// 	fmt.Printf("[%s %s] Message = %s \n", record.EventSource, snsRecord.Timestamp, snsRecord.Message)
	// 	idstr = snsRecord.Message
	// }

	idstr = snsEvent.Records[0].SNS.Message
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
	proj := expression.NamesList(expression.Name("id"), expression.Name("status"), expression.Name("text"), expression.Name("url"), expression.Name("voice"))
	filt := expression.Name("id").Equal(expression.Value(idstr))
	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
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

	// fmt.Println("result.Items = %s", result.Items)

	// err = dynamodbattribute.UnmarshalMap(result.Items, allitems)

	// fmt.Println("allItems = %s", allitems)

	item := Item{}
	resultItem := result.Items[0]
	err = dynamodbattribute.UnmarshalMap(resultItem, &item)

	if err != nil {
		log.Fatalf("Got error unmarshalling: %s", err)
	}

	fmt.Println("Id: ", item.Id)
	fmt.Println("Status:", item.Status)
	fmt.Println("Text:", item.Text)
	fmt.Println("Voice:", item.Voice)
	fmt.Println("URL:", item.Url)
	fmt.Println()

	// Polly
	pollysvc := polly.New(sess)
	input := &polly.SynthesizeSpeechInput{
		OutputFormat: aws.String("mp3"),
		Text:         aws.String(item.Text),
		VoiceId:      aws.String(item.Voice),
	}

	output, err := pollysvc.SynthesizeSpeech(input)
	if err != nil {
		fmt.Println("pollysvc.SynthesizeSpeech error")
		fmt.Print(err.Error())
		os.Exit(1)
	}

	filename := item.Id
	mp3File := filename + ".mp3"
	outFile, err := os.Create("/tmp/" + mp3File)

	if err != nil {
		fmt.Println("Got error creating " + mp3File + ":")
		fmt.Print(err.Error())
		os.Exit(1)
	}

	_, err = io.Copy(outFile, output.AudioStream)
	if err != nil {
		fmt.Println("Got error saving MP3:")
		fmt.Print(err.Error())
		os.Exit(1)
	}
	// upload s3
	uploader := s3manager.NewUploader(sess)
	storeBucket := "jackpollywebsite4"
	s3id := item.Id

	resultS3, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(storeBucket),
		Key:    aws.String(s3id),
		Body:   outFile,
	})

	defer outFile.Close()
	if err != nil {
		fmt.Printf("failed to upload file, %v", err)
	}

	// fmt.Printf("file uploaded to, %s\n", aws.StringValue(resultS3.Location))
	fmt.Printf("file uploaded to, %s\n", resultS3.Location)
}

func main() {
	lambda.Start(handler)

}
