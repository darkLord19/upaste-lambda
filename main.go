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
	Name      string    `json:"name"`
	Data      string    `json:"data"`
	CreatedAt time.Time `json:"created_at"`
}

var db *dynamodb.DynamoDB
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
	if req.Headers["Content-Type"] != "application/json" || req.Headers["content-type"] != "application/json" {
		return clientError(http.StatusNotAcceptable)
	}

	var paste Paste
	err := json.Unmarshal([]byte(req.Body), &paste)
	if err != nil {
		return clientError(http.StatusBadRequest)
	}
	ct := time.Now()
	item := Paste{
		Name:      "upaste_" + fmt.Sprint(ct.Unix()),
		Data:      paste.Data,
		CreatedAt: ct,
	}
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return serverError(err)
	}

	tableName := "pastes"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = db.PutItem(input)
	if err != nil {
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
	region := "ap-south-1"
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return
	}
	db = dynamodb.New(awsSession)
	lambda.Start(router)
}
