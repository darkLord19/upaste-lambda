package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

func serverError(err error) (events.APIGatewayProxyResponse, error) {
	log.Println("ERROR ", err.Error())

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       err.Error(),
	}, nil
}

func clientError(err error) (events.APIGatewayProxyResponse, error) {
	log.Println("ERROR ", err.Error())

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Body:       err.Error(),
	}, nil
}

func createPaste(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("ERROR: ", "HEADERS: ", req.Headers, "BODY:", req.Body)
	if req.Headers["Content-Type"] != "application/json" && req.Headers["content-type"] != "application/json" {
		return clientError(fmt.Errorf("%s", "Bad Request"))
	}

	var paste Paste
	err := json.Unmarshal([]byte(req.Body), &paste)
	if err != nil {
		return clientError(fmt.Errorf("%s %s", err.Error(), req.Body))
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

func main() {
	region := "ap-south-1"
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return
	}
	db = dynamodb.New(awsSession)
	lambda.Start(createPaste)
}
