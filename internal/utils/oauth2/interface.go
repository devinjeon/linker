package oauth2

// Token is struct having fields, accessToken and refreshToken.
type Token struct {
	AccessToken  string `dynamodbav:"access_token"`
	RefreshToken string `dynamodbav:"refresh_token"`
}

type oauth2 interface {
	GetAuthorizeURI() string
	ExchangeToken(code string) (*Token, error)
	ValidateToken(token *Token) (bool, error)
	RefreshToken(token *Token) (bool, error)
	RevokeToken(token *Token) (bool, error)
	UserEmail(token *Token) (string, error)
}
