package auth

import (
	"github.com/aws/aws-lambda-go/events"
	"linker/utils/oauth2"
	"os"
)

var o oauth2.GitHub

// Response is of type APIGatewayProxyResponse
type Response = events.APIGatewayProxyResponse

// Request is of type APIGatewayProxyRequest
type Request = events.APIGatewayProxyRequest

func init() {
	clientID := os.Getenv("OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("OAUTH_CLIENT_SECRET")
	redirectURI := os.Getenv("OAUTH_REDIRECT_URI")

	o = oauth2.GitHub{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
	}
}

// Handler returns links API response
func Handler(req Request) (Response, error) {
	path := req.PathParameters["proxy"]

	switch path {
	case "signin":
		return signin()
	case "exchange":
		return exchange(req)
	default:
		return Response{StatusCode: 400}, nil
	}
}

func signin() (Response, error) {
	resp := Response{
		StatusCode: 301,
		Headers: map[string]string{
			"Location": o.GetAuthorizeURI(),
		},
	}

	return resp, nil
}

func exchange(req Request) (Response, error) {
	code, ok := req.QueryStringParameters["code"]
	if !ok {
		return Response{StatusCode: 400}, nil
	}
	token, err := o.ExchangeToken(code)
	if err != nil {
		return Response{StatusCode: 500}, nil
	}

	return Response{
		Body:       token.AccessToken,
		StatusCode: 200,
	}, nil
}
