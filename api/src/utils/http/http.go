package http

import (
	"bytes"
	"fmt"
	"net/http"
	u "net/url"
	"runtime/debug"
)

var client = &http.Client{}

const (
	// HTTPGet represents GET method.
	HTTPGet = "GET"
	// HTTPPost represents POST method.
	HTTPPost = "POST"
	// HTTPDelete represents DELETE method.
	HTTPDelete = "DELETE"
	// HTTPPut represents PUT method.
	HTTPPut = "PUT"
)

// Response is result type that HTTP request returns.
type Response struct {
	StatusCode int
	Body       []byte
}

func request(method string, url string, queryParams map[string]string, headers map[string]string, data []byte) (response *Response, err error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("[ERROR]http.%s(%s, %v, %v):\n%v\n", method, url, headers, queryParams, err)
			debug.PrintStack()
			response = nil
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

	req, err := http.NewRequest(method, newURL.String(), bytes.NewBuffer(data))
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

	response = &Response{
		Body:       buf.Bytes(),
		StatusCode: resp.StatusCode,
	}
	return
}

// Get requests HTTP call with GET method
func Get(url string, queryParams map[string]string, headers map[string]string, body []byte) (response *Response, err error) {
	return request(HTTPGet, url, queryParams, headers, body)
}

// Post requests HTTP call with POST method
func Post(url string, queryParams map[string]string, headers map[string]string, body []byte) (response *Response, err error) {
	return request(HTTPGet, url, queryParams, headers, body)
}

// Delete requests HTTP call with DELETE method
func Delete(url string, queryParams map[string]string, headers map[string]string, body []byte) (response *Response, err error) {
	return request(HTTPGet, url, queryParams, headers, body)
}

// Put requests HTTP call with PUT method
func Put(url string, queryParams map[string]string, headers map[string]string, body []byte) (response *Response, err error) {
	return request(HTTPGet, url, queryParams, headers, body)
}
