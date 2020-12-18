package links

import (
	"encoding/json"
	db "linker/utils/dynamodb"
	"strings"

	"github.com/aws/aws-lambda-go/events"
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

type newLink struct {
	URL string `json:"url"`
	TTL int    `json:"ttl"`
}

func redirect(req Request) (Response, error) {
	id := strings.TrimPrefix(req.Path, "/")

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
	id := strings.TrimPrefix(req.Path, "/")
	body := req.Body

	var data newLink
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return Response{StatusCode: 400}, err
	}

	err := putURL(id, data.URL, data.TTL)
	switch err {
	case db.ErrDBOperation:
		return Response{StatusCode: 500}, err
	case db.ErrMarshalling:
		return Response{StatusCode: 500}, err
	}

	return Response{StatusCode: 204}, nil
}

func delete(req Request) (Response, error) {
	id := strings.TrimPrefix(req.Path, "/")

	err := deleteURL(id)
	switch err {
	case db.ErrDBOperation:
		return Response{StatusCode: 500}, err
	case db.ErrNotFoundItem:
		return Response{StatusCode: 404}, nil
	case db.ErrUnmarshalling:
		return Response{StatusCode: 500}, err
	}

	return Response{StatusCode: 204}, nil
}
