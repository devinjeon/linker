package oauth2

import (
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
