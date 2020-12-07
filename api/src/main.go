package main

import (
	"encoding/json"
	"os"

	"linker/db"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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

func upsert(req Request) (Response, error) {
	id := req.PathParameters["id"]
	body := req.Body

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return Response{StatusCode: 400}, err
	}
	url := data["url"].(string)

	err := dbClient.PutURL(id, url)
	switch err {
	case db.ErrDBOperation:
		return Response{StatusCode: 500}, err
	case db.ErrMarshalling:
		return Response{StatusCode: 500}, err
	}

	return Response{StatusCode: 200}, nil
}

func delete(req Request) (Response, error) {
	id := req.PathParameters["id"]

	err := dbClient.DeleteURL(id)
	switch err {
	case db.ErrDBOperation:
		return Response{StatusCode: 500}, err
	case db.ErrNotFoundItem:
		return Response{StatusCode: 404}, nil
	case db.ErrUnmarshalling:
		return Response{StatusCode: 500}, err
	}

	return Response{StatusCode: 200}, nil
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(req Request) (Response, error) {
	method := req.HTTPMethod
	switch method {
	case "GET":
		return redirect(req)
	case "POST":
		return upsert(req)
	case "DELETE":
		return delete(req)
	case "PUT":
		return upsert(req)
	default:
		return Response{StatusCode: 405}, nil
	}
}

func main() {
	lambda.Start(Handler)
}
