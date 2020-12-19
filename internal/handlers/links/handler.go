package links

import (
	"encoding/json"
	"strings"

	m "github.com/devinjeon/linker/internal/middleware"
	db "github.com/devinjeon/linker/internal/utils/dynamodb"
)

type (
	response = m.Response
	request  = m.Request
)

// Handler returns links API response
func Handler(req request) (response, error) {
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
		return response{StatusCode: 405}, nil
	}
}

type newLink struct {
	URL string `json:"url"`
	TTL int    `json:"ttl"`
}

func redirect(req request) (response, error) {
	id := strings.TrimPrefix(req.Path, "/")

	url, err := getURL(id)
	switch err {
	case db.ErrDBOperation:
		return response{StatusCode: 500}, err
	case db.ErrNotFoundItem:
		return response{StatusCode: 404}, nil
	case db.ErrUnmarshalling:
		return response{StatusCode: 500}, err
	}

	resp := response{
		StatusCode: 301,
		Headers: map[string]string{
			"Location": url,
		},
	}

	return resp, nil
}

func upsert(req request) (response, error) {
	id := strings.TrimPrefix(req.Path, "/")
	body := req.Body

	var data newLink
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return response{StatusCode: 400}, err
	}

	err := putURL(id, data.URL, data.TTL)
	switch err {
	case db.ErrDBOperation:
		return response{StatusCode: 500}, err
	case db.ErrMarshalling:
		return response{StatusCode: 500}, err
	}

	return response{StatusCode: 204}, nil
}

func delete(req request) (response, error) {
	id := strings.TrimPrefix(req.Path, "/")

	err := deleteURL(id)
	switch err {
	case db.ErrDBOperation:
		return response{StatusCode: 500}, err
	case db.ErrNotFoundItem:
		return response{StatusCode: 404}, nil
	case db.ErrUnmarshalling:
		return response{StatusCode: 500}, err
	}

	return response{StatusCode: 204}, nil
}
