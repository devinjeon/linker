package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"os"

	db "github.com/devinjeon/linker/internal/utils/dynamodb"
	"github.com/devinjeon/linker/internal/utils/oauth2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var sessionTable = os.Getenv("DYNAMODB_SESSION_TABLE_NAME")
var clientID = os.Getenv("OAUTH_CLIENT_ID")
var clientSecret = os.Getenv("OAUTH_CLIENT_SECRET")
var redirectURI = os.Getenv("OAUTH_REDIRECT_URI")

// OAuth2 is OAuth2 client to handle OAuth2 flow
var OAuth2 = oauth2.GitHub{
	ClientID:     clientID,
	ClientSecret: clientSecret,
	RedirectURI:  redirectURI,
}

var c = db.NewDB(sessionTable)

// Session is struct to handle session
type Session struct {
	ID    string       `dynamodbav:"session_id"`
	Token oauth2.Token `dynamodbav:"token"`
}

// UserEmail gets user email from access token
func (s *Session) UserEmail() (string, error) {
	return OAuth2.UserEmail(s.Token)
}

// GetSession gets session from id
func getSession(id string) (*Session, error) {
	key := map[string]*dynamodb.AttributeValue{
		"session_id": {
			S: aws.String(id),
		},
	}
	sess := &Session{}
	err := c.GetItem(key, &sess)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func validateSession(session Session) (bool, error) {
	return OAuth2.ValidateToken(session.Token)
}

// NewSession creates a new session
func NewSession(token *oauth2.Token) (*Session, error) {
	sessID, err := generateSessionID()
	if err != nil {
		return nil, err
	}

	sess := &Session{
		ID:    sessID,
		Token: *token,
	}

	err = c.PutItem(*sess)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func generateSessionID() (string, error) {
	// generate session ID
	length := 64
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	sessID := base64.URLEncoding.EncodeToString(b)
	return sessID, err
}

// RemoveSession removes session from DB and revokes Token
func RemoveSession(session Session) error {
	key := map[string]*dynamodb.AttributeValue{
		"session_id": {
			S: aws.String(session.ID),
		},
	}

	OAuth2.RevokeToken(session.Token)
	err := c.DeleteItem(key)
	if err != nil {
		return db.ErrDBOperation
	}

	return nil
}
