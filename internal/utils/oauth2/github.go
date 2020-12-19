package oauth2

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	u "net/url"

	"github.com/devinjeon/linker/internal/utils/http"
)

// GitHub is struct to handle GitHub OAuth2 flow
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
	if o.RedirectURI != "" {
		q.Set("redirect_uri", o.RedirectURI)
	}
	authURI.RawQuery = q.Encode()
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
		return nil, err
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
		RefreshToken: "null",
	}
	return token, nil
}

// tokenAPI returns GitHub OAuth2 token API URL
func (o *GitHub) tokenAPI() string {
	return fmt.Sprintf("https://api.github.com/applications/%s/token", o.ClientID)
}

// ValidateToken checks if the accessToken is valid.
func (o *GitHub) ValidateToken(token Token) (bool, error) {
	apiURL := o.tokenAPI()

	auth := o.ClientID + ":" + o.ClientSecret
	basicAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	headers := map[string]string{
		"Authorization": "Basic " + basicAuth,
	}
	data := map[string]string{
		"access_token": token.AccessToken,
	}

	dataJSON, _ := json.Marshal(data)
	resp, err := http.Post(apiURL, nil, headers, dataJSON)
	if err != nil {
		return false, err
	}

	if resp.StatusCode == 200 {
		return true, nil
	}
	if resp.StatusCode == 404 {
		return false, nil
	}

	return false, ErrOAuthServer
}

// RefreshToken does nothing because GitHub OAuth2 doesn't expire access token.
func (o *GitHub) RefreshToken(token Token) (bool, error) {
	return false, nil
}

// RevokeToken revokes access token and authorization
func (o *GitHub) RevokeToken(token Token) (bool, error) {
	apiURL := o.tokenAPI()

	auth := o.ClientID + ":" + o.ClientSecret
	basicAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	headers := map[string]string{
		"Authorization": "Basic " + basicAuth,
	}
	data := map[string]string{
		"access_token": token.AccessToken,
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return false, err
	}

	resp, err := http.Delete(apiURL, nil, headers, dataJSON)
	if err != nil {
		return false, err
	}

	// Success
	if resp.StatusCode == 204 {
		return true, nil
	}

	return false, ErrOAuthServer
}

// UserEmail returns user email by call user API with token
func (o *GitHub) UserEmail(token Token) (string, error) {
	apiURL := "https://api.github.com/user"

	headers := map[string]string{
		"Authorization": fmt.Sprintf("bearer %s", token.AccessToken),
	}

	resp, err := http.Get(apiURL, nil, headers, nil)
	if err != nil {
		return "", err
	}

	if resp.StatusCode == 401 {
		return "", ErrUnauthorized
	}
	if resp.StatusCode != 200 {
		return "", ErrOAuthServer
	}

	// Parse reponse body
	var responseData map[string]interface{}
	if err := json.Unmarshal(resp.Body, &responseData); err != nil {
		return "", err
	}

	email, isValid := responseData["email"]
	if !isValid {
		errorMessage, ok := responseData["message"]
		if !ok {
			return "", ErrOAuthServer
		}
		err = fmt.Errorf("statusCode: %d, %s", resp.StatusCode, errorMessage.(string))
		return "", err
	}

	return email.(string), nil
}
