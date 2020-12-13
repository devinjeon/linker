package oauth2

import (
	"encoding/json"
	"fmt"
	"linker/utils/http"
	u "net/url"
	"runtime/debug"
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
	var responseData map[string]interface{}
	if err := json.Unmarshal(resp.Body, &responseData); err != nil {
		return nil, ErrMarshalling
	}

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("[ERROR]oauth2.github:\nresponse: %s\nerror: %s\n", string(resp.Body), err)
			debug.PrintStack()
		}
	}()

	token = &Token{
		AccessToken:  responseData["access_token"].(string),
		RefreshToken: "",
		TokenType:    responseData["token_type"].(string),
	}
	return token, nil
}
