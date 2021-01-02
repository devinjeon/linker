package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

const loginSessionTime = 3600 * 24 * 90

func init() {
	gob.Register(oauth2.Token{})
}

// New creates handlers for authentication
func New(redirectURL string, oauth2 oauth2.Config) Handlers {
	handlers := Handlers{
		redirectURL: redirectURL,
		oauth2:      oauth2,
	}

	return handlers
}

// Handlers is struct including handler methods.
type Handlers struct {
	redirectURL string
	oauth2      oauth2.Config
}

// User is map[string]interface{} type storing user information
type User map[string]interface{}

// Auth stores user information
type Auth struct {
	// User is user ID of OAuth2 service
	User string
	// Token is access token
	Token oauth2.Token
}

func (h *Handlers) getSession(c *gin.Context) sessions.Session {
	session := sessions.Default(c)
	return session
}

func generateState() (string, error) {
	length := 32
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.URLEncoding.EncodeToString(b)
	return state, nil
}

// SignIn is a handler to sign in.
func (h *Handlers) SignIn(c *gin.Context) {
	session := h.getSession(c)
	state, err := generateState()
	if err != nil {
		c.AbortWithError(500, err)
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("state", state, 0, "/", "", false, true)
	c.SetCookie("referrer", c.Request.Referer(), 0, "/", "", false, true)
	if err := session.Save(); err != nil {
		c.AbortWithError(500, err)
	}

	authorizeURL := h.oauth2.AuthCodeURL(state, oauth2.AccessTypeOnline)

	c.Redirect(302, authorizeURL)
}

// SignOut is a handler to sign out.
func (h *Handlers) SignOut(c *gin.Context) {
	c.SetCookie("is_logged_in", "false", -1, "/", "", false, false)
	session := h.getSession(c)
	session.Clear()
	if err := session.Save(); err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.Redirect(302, h.redirectURL)
	return
}

func (h *Handlers) userInfo(token *oauth2.Token) (User, error) {
	apiURL := "https://api.github.com/user"

	client := h.oauth2.Client(oauth2.NoContext, token)
	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Parse reponse body
	var user User

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(string(body))
	}

	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}

	return user, nil
}

// Exchange is a handler to exchange token between OAuth2 service and Linker.
func (h *Handlers) Exchange(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			c.String(500, "Error: %v", r.(error))
			return
		}
	}()

	// Get and Delete cookie values
	state, _ := c.Cookie("state")
	referrer, _ := c.Cookie("referrer")
	c.SetCookie("state", "", -1, "", "", false, true)
	c.SetCookie("referrer", "", -1, "", "", false, true)

	// Check code from OAuth2 service
	code := c.Query("code")
	if code == "" {
		c.AbortWithStatus(400)
		return
	}

	// Check state
	if state != c.Query("state") {
		c.AbortWithStatus(400)
		return
	}

	// Exchange token
	token, err := h.oauth2.Exchange(oauth2.NoContext, code)
	if err != nil {
		c.String(500, "Failed to exchange token: %v", err)
		return
	}

	// Get user info
	userInfo, err := h.userInfo(token)
	if err != nil {
		c.String(500, "Failed to get user info: %v", err)
		return
	}

	// Get userID
	userID, ok := userInfo["id"]
	if !ok {
		c.String(500, "Failed to get userID")
		return
	}

	// Encode userID
	user := base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%v", userID)))

	// Save user and token to session
	session := h.getSession(c)
	session.Clear()
	session.Set("user", user)
	session.Set("token", token)
	if err := session.Save(); err != nil {
		c.String(500, "Failed to save session: %v", err)
		return
	}

	// Set login status
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("is_logged_in", "true", loginSessionTime, "/", "", false, false)

	// Get redirect URL
	var redirectURL string
	if referrer != "" {
		redirectURL = referrer
	} else {
		redirectURL = h.redirectURL
	}

	c.Redirect(302, redirectURL)
}

// RequireAuth is a middleware to check if the current user is authenticated
func (h *Handlers) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := h.getSession(c)

		isValid := true
		user, ok := session.Get("user").(string)
		if !ok || user == "" {
			isValid = false
		}
		token, ok := session.Get("token").(oauth2.Token)
		if !ok {
			isValid = false
		}

		if !isValid {
			session.Clear()
			session.Save()
			c.SetCookie("is_logged_in", "false", -1, "/", "", false, false)
			c.AbortWithStatus(401)
			return
		}

		auth := Auth{
			User:  user,
			Token: token,
		}
		c.Set("auth", auth)
	}
}
