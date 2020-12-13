package oauth2

import "errors"

var (
	// ErrOAuthServer is a error that represents "OAuth server error"
	ErrOAuthServer = errors.New("oauth2: OAuth server error")
	// ErrMarshalling is a error that represents "JSON Marshalling error"
	ErrMarshalling = errors.New("oauth2: JSON Marshalling error")
)
