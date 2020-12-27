package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"os"

	db "github.com/devinjeon/linker/internal/utils/dynamodb"
	"github.com/devinjeon/linker/internal/utils/oauth2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/gin-gonic/gin"
)

var sessionTable = os.Getenv("DYNAMODB_SESSION_TABLE_NAME")
var clientID = os.Getenv("OAUTH_CLIENT_ID")
var clientSecret = os.Getenv("OAUTH_CLIENT_SECRET")
var redirectURI = os.Getenv("OAUTH_REDIRECT_URI")

// OAuth2 is OAuth2 client to handle OAuth2 flow.
var OAuth2 = oauth2.GitHub{
	ClientID:     clientID,
	ClientSecret: clientSecret,
	RedirectURI:  redirectURI,
}

var sessions = db.NewDB(sessionTable)

// Session is struct to handle session.
type Session struct {
	ID    string       `dynamodbav:"session_id"`
	Token oauth2.Token `dynamodbav:"token"`
}

// UserEmail gets user email from access token.
func (s *Session) UserEmail() (string, error) {
	return OAuth2.UserEmail(s.Token)
}

func loadSessionByID(id string) (*Session, error) {
	key := map[string]*dynamodb.AttributeValue{
		"session_id": {
			S: aws.String(id),
		},
	}
	sess := &Session{}
	err := sessions.GetItem(key, &sess)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func validateSession(session Session) (bool, error) {
	return OAuth2.ValidateToken(session.Token)
}

// GetSession gets Session from context.
func GetSession(c *gin.Context) (*Session, bool) {
	s, ok := c.Get("session")
	if !ok {
		return nil, false
	}

	return s.(*Session), true
}

// NewSession creates a new session.
func NewSession(token *oauth2.Token, c *gin.Context) error {
	sessID, err := generateSessionID()
	if err != nil {
		return err
	}

	sess := &Session{
		ID:    sessID,
		Token: *token,
	}

	err = sessions.PutItem(*sess)
	if err != nil {
		return err
	}

	c.Set("session", *sess)
	SetCookie("session_id", sess.ID, c)

	return nil
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

// RemoveSession removes session from DB and revokes Token.
func RemoveSession(c *gin.Context) error {
	UnsetCookie("session_id", c)

	s, ok := c.Get("session")
	if !ok {
		return errors.New("No session")
	}
	sess := s.(Session)
	key := map[string]*dynamodb.AttributeValue{
		"session_id": {
			S: aws.String(sess.ID),
		},
	}

	OAuth2.RevokeToken(sess.Token)
	err := sessions.DeleteItem(key)
	if err != nil {
		return db.ErrDBOperation
	}

	return nil
}

// RequireSession is wrapper for handlers that must have valid session.
func RequireSession(f func(c *gin.Context)) func(c *gin.Context) {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err == http.ErrNoCookie {
			c.Status(400)
			return
		}

		sess, err := loadSessionByID(sessionID)
		if err != nil {
			RemoveSession(c)
			c.Status(401)
			return
		}

		if isValid, _ := validateSession(*sess); !isValid {
			RemoveSession(c)
			c.Status(401)
			return
		}

		c.Set("session", sess)
		f(c)
	}
}
