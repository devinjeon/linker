package main

import (
	"strings"

	"linker/apis/auth"
	"linker/apis/links"
	m "linker/middleware"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// This Lambda function is triggered by API Gateway
// so it should use Request and Response types from events.APIGatewayProxy*.

type (
	response = m.Response
	request  = m.Request
	// Handler is type of handler function
	Handler func(request) (response, error)
)

type router struct {
	handlers map[string]Handler
}

func newRouter() router {
	r := router{
		handlers: make(map[string]Handler),
	}
	return r
}

func (r *router) addHandler(path string, handler Handler) {
	if path == "" {
		path = "/"
	}

	if path[0] != '/' {
		path = "/" + path
	}

	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	r.handlers[path] = handler
}

func (r *router) route(orgReq events.APIGatewayProxyRequest) (response, error) {
	req := m.WrapAPIGatewayProxyRequest(orgReq)

	badRequest := response{StatusCode: 400}
	path := req.Path
	if path == "" {
		path = "/"
	}
	if path[0] != '/' {
		path = "/" + path
	}

	for p, h := range r.handlers {
		if !strings.HasPrefix(path, p) {
			continue
		}

		// NOTE: request path should not match the case
		//       that handler path is "/{req.path}xxxx/yyyy"
		if path[len(path)-1] != '/' {
			path = path + "/"
		}
		if strings.HasPrefix(path, p+"/") {
			// NOTE: handler gets only subpath
			subpath := strings.TrimPrefix(path, p)
			subpath = strings.TrimSuffix(subpath, "/")
			req.Path = subpath

			return h(req)
		}
	}

	return badRequest, nil
}

func main() {
	r := newRouter()

	r.addHandler("/links", links.Handler)
	r.addHandler("/auth", auth.Handler)

	lambda.Start(r.route)
}
