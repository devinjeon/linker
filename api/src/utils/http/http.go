package main

import (
	"bytes"
	"fmt"
	"net/http"
	u "net/url"
)

var client = &http.Client{}

func get(url string, headers map[string]string, queryParams map[string]string) (responseBody string, statusCode int) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("[ERROR]http.Get(%s, %v, %v):\n%s\n", url, headers, queryParams, err)
			responseBody = ""
			statusCode = -1
		}
	}()

	newURL, err := u.Parse(url)
	if err != nil {
		panic(err)
	}

	q := newURL.Query()
	if q != nil {
		for k, v := range queryParams {
			q.Set(k, v)
		}
	}
	newURL.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", newURL.String(), nil)
	if err != nil {
		panic(err)
	}

	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		panic(err)
	}

	responseBody = buf.String()
	statusCode = resp.StatusCode
	return
}
