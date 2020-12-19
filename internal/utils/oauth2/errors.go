package oauth2

import "errors"

var (
	// ErrOAuthServer is a error that represents "OAuth server error"
	ErrOAuthServer = errors.New("oauth2: OAuth server error")
	// ErrUnauthorized is a error that represents "Unauthorized"
	ErrUnauthorized = errors.New("oauth2: Unauthorized")
)
