package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// This Lambda function is triggered by API Gateway
// so it should use Request and Response types from events.APIGatewayProxy*.

// Response is of type APIGatewayProxyResponse
type Response events.APIGatewayProxyResponse

// Request is of type APIGatewayProxyRequest
type Request events.APIGatewayProxyRequest

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))

var db = dynamodb.New(sess)
var tableName = os.Getenv("DYNAMODB_TABLE_NAME")

// Link struct is a model for items from DynamoDB table
type Link struct {
	ID  string
	URL string
}

func redirect(req Request) (Response, error) {
	id := req.PathParameters["id"]

	result, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})

	resp := Response{
		StatusCode: 404,
	}
	if err != nil {
		fmt.Println(err.Error())
		return Response{StatusCode: 500}, err
	}

	if result.Item == nil {
		fmt.Println("Could not find '" + id + "'")
		return Response{StatusCode: 404}, nil
	}

	item := Link{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		fmt.Printf("Failed to unmarshal Record, %v", err)
		return Response{StatusCode: 500}, err
	}

	resp = Response{
		StatusCode: 301,
		Headers: map[string]string{
			"Location": item.URL,
		},
	}

	return resp, nil

}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(req Request) (Response, error) {
	method := req.HTTPMethod
	switch method {
	case "GET":
		return redirect(req)
	case "POST":
		return redirect(req)
	case "DELETE":
		return redirect(req)
	case "PUT":
		return redirect(req)
	default:
		return Response{StatusCode: 405}, nil
	}
}

func main() {
	lambda.Start(Handler)
}
