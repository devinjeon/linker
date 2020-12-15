package links

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	db "linker/utils/dynamodb"
)

// Response is of type APIGatewayProxyResponse
type Response = events.APIGatewayProxyResponse

// Request is of type APIGatewayProxyRequest
type Request = events.APIGatewayProxyRequest

// Handler returns links API response
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

func redirect(req Request) (Response, error) {
	id := req.PathParameters["proxy"]

	url, err := getURL(id)
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
	id := req.PathParameters["proxy"]
	body := req.Body

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return Response{StatusCode: 400}, err
	}
	url := data["url"].(string)

	err := putURL(id, url)
	switch err {
	case db.ErrDBOperation:
		return Response{StatusCode: 500}, err
	case db.ErrMarshalling:
		return Response{StatusCode: 500}, err
	}

	return Response{StatusCode: 200}, nil
}

func delete(req Request) (Response, error) {
	id := req.PathParameters["proxy"]

	err := deleteURL(id)
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