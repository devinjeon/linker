package middleware

import "github.com/aws/aws-lambda-go/events"

// Response is alias of events.APIGatewayProxyResponse
type Response = events.APIGatewayProxyResponse

// Request is extended type of APIGatewayProxyRequest
type Request struct {
	events.APIGatewayProxyRequest
	Session *Session
	Cookies map[string]string
}

// WrapAPIGatewayProxyRequest converts APIGatewayProxyRequest to middleware.Request
func WrapAPIGatewayProxyRequest(req events.APIGatewayProxyRequest) Request {
	cookies := make(map[string]string)
	cookieHeader, ok := req.Headers["cookie"]
	if ok {
		cookies = parseCookie(cookieHeader)
	}

	sessionID, _ := cookies["session_id"]
	sess, _ := getSession(sessionID)
	return Request{
		APIGatewayProxyRequest: req,
		Session:                sess,
		Cookies:                cookies,
	}
}
