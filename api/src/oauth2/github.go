package oauth2

import (
	"encoding/json"
	"errors"
	"fmt"
	"linker/utils/http"
	u "net/url"
)

// GitHub is struct for handling GitHub OAuth2 flow
type GitHub struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

// GetAuthorizeURI returns URI to authorize
func (o *GitHub) GetAuthorizeURI() string {
	baseURL := "https://github.com/login/oauth/authorize"
	authURI, _ := u.Parse(baseURL)

	q := authURI.Query()
	q.Set("client_id", o.ClientID)
	authURI.RawQuery = q.Encode()
	if o.RedirectURI != "" {
		q.Set("redirect_uri", o.RedirectURI)
	}
	return authURI.String()
}

// ExchangeToken acquires access token
func (o *GitHub) ExchangeToken(code string) (token *Token, err error) {
	apiURL := "https://github.com/login/oauth/access_token"

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}
	data := map[string]string{
		"client_id":     o.ClientID,
		"client_secret": o.ClientSecret,
		"code":          code,
	}

	dataJSON, _ := json.Marshal(data)
	resp, err := http.Post(apiURL, nil, headers, dataJSON)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, ErrOAuthServer
	}

	// Parse reponse body
	var responseData map[string]interface{}
	if err := json.Unmarshal(resp.Body, &responseData); err != nil {
		return nil, ErrMarshalling
	}

	if _, isValid := responseData["access_token"]; !isValid {
		errorMessage, _ := responseData["error_description"]
		if errorMessage == nil {
			return nil, ErrOAuthServer
		}
		return nil, errors.New(errorMessage.(string))
	}

	token = &Token{
		AccessToken:  responseData["access_token"].(string),
		RefreshToken: "",
		TokenType:    responseData["token_type"].(string),
	}
	return token, nil
}
