package auth

import (
	"fmt"
	"net/http"
	"os"

	m "github.com/devinjeon/linker/internal/middleware"
	"github.com/gin-gonic/gin"
)

var linkerDomain = os.Getenv("LINKER_DOMAIN")
var linkerURL = fmt.Sprintf("https://%s", linkerDomain)

// SignIn is a handler to sign in.
func SignIn(c *gin.Context) {
	c.Status(301)
	c.Header("Cache-control", "no-cache")
	_, err := c.Cookie("session_id")
	if err != http.ErrNoCookie {
		c.Header("Location", linkerURL)
		return
	}
	c.Header("Location", m.OAuth2.GetAuthorizeURI())
	return
}

// SignOut is a handler to sign out.
func SignOut(c *gin.Context) {
	m.RemoveSession(c)
	c.Header("Cache-control", "no-cache")
	c.Status(204)
}

// Exchange is a handler to exchange token between OAuth2 service and Linker.
func Exchange(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.Status(400)
		return
	}

	token, err := m.OAuth2.ExchangeToken(code)
	if err != nil {
		c.Status(500)
		return
	}

	err = m.NewSession(token, c)
	if err != nil {
		c.Status(500)
		return
	}

	c.Status(301)
	c.Header("Location", linkerURL)
	c.Header("Cache-control", "no-cache")
}
