package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
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
var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

func serverError(err error) (events.APIGatewayProxyResponse, error) {
	errorLogger.Println(err.Error())

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func createPaste(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if req.Headers["Content-Type"] != "application/json" {
		return clientError(http.StatusNotAcceptable)
	}

	var paste Paste
	err := json.Unmarshal([]byte(req.Body), paste)
	ct := time.Now()
	item := Paste{
		Name:      "upaste_" + fmt.Sprint(ct.Unix()),
		Data:      paste.Data,
		CreatedAt: ct,
	}
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Fatalf("Got error marshalling new movie item: %s", err)
		return serverError(err)
	}

	tableName := "pastes"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = db.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
		return serverError(err)
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       http.StatusText(http.StatusCreated),
	}, nil
}

func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "POST":
		return createPaste(req)
	default:
		return clientError(http.StatusMethodNotAllowed)
	}
}

func main() {
	lambda.Start(router)
}
