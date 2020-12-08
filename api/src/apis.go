package main

import (
	"encoding/json"
	"linker/db"
	"os"
)

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
