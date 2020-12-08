package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// This Lambda function is triggered by API Gateway
// so it should use Request and Response types from events.APIGatewayProxy*.

// Response is of type APIGatewayProxyResponse
type Response events.APIGatewayProxyResponse

// Request is of type APIGatewayProxyRequest
type Request events.APIGatewayProxyRequest

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
