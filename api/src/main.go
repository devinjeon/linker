package main

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"linker/db"
)

// This Lambda function is triggered by API Gateway
// so it should use Request and Response types from events.APIGatewayProxy*.

// Response is of type APIGatewayProxyResponse
type Response events.APIGatewayProxyResponse

// Request is of type APIGatewayProxyRequest
type Request events.APIGatewayProxyRequest

var tableName = os.Getenv("DYNAMODB_TABLE_NAME")
var dbClient = db.GetDB(tableName)

func redirect(req Request) (Response, error) {
	id := req.PathParameters["id"]

	url, err := dbClient.GetURL(id)
	switch err {
	case db.ErrDBOperation:
		return Response{StatusCode: 500}, err
	case db.ErrNotFoundItem:
		return Response{StatusCode: 404}, nil
	case db.ErrUnmarshalling:
		return Response{StatusCode: 500}, err
	}

	resp := Response{
		StatusCode: 301,
		Headers: map[string]string{
			"Location": url,
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
		return Response{StatusCode: 405}, nil
	case "DELETE":
		return Response{StatusCode: 405}, nil
	case "PUT":
		return Response{StatusCode: 405}, nil
	default:
		return Response{StatusCode: 405}, nil
	}
}

func main() {
	lambda.Start(Handler)
}
