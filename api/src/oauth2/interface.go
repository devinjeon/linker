package oauth2

// Token is struct having fields, accessToken and refreshToken
type Token struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
}

type oauth2 interface {
	GetAuthorizeURI() string
	ExchangeToken(code string) (*Token, error)
	ValidateToken(accessToken string) (bool, error)
	RefreshToken(refreshToken string) (*Token, error)
	RevokeToken(accessToken string) bool
}
