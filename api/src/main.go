package main

import (
	"linker/apis/links"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// This Lambda function is triggered by API Gateway
// so it should use Request and Response types from events.APIGatewayProxy*.

// Response is of type APIGatewayProxyResponse
type Response = events.APIGatewayProxyResponse

// Request is of type APIGatewayProxyRequest
type Request = events.APIGatewayProxyRequest

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(req Request) (Response, error) {
	badRequest := Response{StatusCode: 400}

	path, ok := req.PathParameters["proxy"]
	if !ok {
		return badRequest, nil
	}

	pathSplited := strings.SplitN(path, "/", 2)
	if len(pathSplited) == 0 {
		return badRequest, nil
	}
	subPath := ""
	if len(pathSplited) > 1 {
		subPath = pathSplited[1]
	}
	req.PathParameters["proxy"] = subPath

	resource := pathSplited[0]
	switch resource {
	case "links":
		return links.Handler(req)
	default:
		return badRequest, nil
	}
}

func main() {
	lambda.Start(Handler)
}
