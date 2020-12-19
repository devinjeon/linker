package auth

import (
	m "github.com/devinjeon/linker/internal/middleware"
)

type request = m.Request
type response = m.Response

// Handler returns links API response
func Handler(req request) (response, error) {
	path := req.Path

	switch path {
	case "/signin":
		return signin(req)
	case "/exchange":
		return exchange(req)
	default:
		return response{StatusCode: 400}, nil
	}
}

func signin(req request) (response, error) {
	if req.Session != nil {
		return response{
			StatusCode: 301,
			Headers: map[string]string{
				"Location":      m.WebRootURI,
				"Cache-control": "no-cache",
			},
		}, nil
	}

	resp := response{
		StatusCode: 301,
		Headers: map[string]string{
			"Location":      m.OAuth2.GetAuthorizeURI(),
			"Cache-control": "no-cache",
		},
	}

	return resp, nil
}

func exchange(req request) (response, error) {
	code, ok := req.QueryStringParameters["code"]
	if !ok {
		return response{StatusCode: 400}, nil
	}
	token, err := m.OAuth2.ExchangeToken(code)
	if err != nil {
		return response{StatusCode: 500}, err
	}

	sess, err := m.NewSession(token)
	if err != nil {
		return response{StatusCode: 500}, err
	}
	resp := response{
		Body:       token.AccessToken,
		StatusCode: 200,
	}
	m.SetCookie("session_id", sess.ID, &resp)

	return resp, nil
}
