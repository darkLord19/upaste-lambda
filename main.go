package main

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Paste struct {
	Name      string
	Data      string
	CreatedAt time.Time
}

var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("ap-south-1"))

func createPaste(paste string) {
	ct := time.Now()
	item := Paste{
		Name:      "upaste_" + fmt.Sprint(ct.Unix()),
		Data:      paste,
		CreatedAt: ct,
	}
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Fatalf("Got error marshalling new movie item: %s", err)
		return
	}
	// Create item in table Movies
	tableName := "pastes"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = db.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
		return
	}
}

func main() {
	lambda.Start(createPaste)
}
