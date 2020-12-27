package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var domain = os.Getenv("LINKER_DOMAIN")

const maxAge int = 3600 * 24 * 90

// SetCookie appends 'Set-Cookie' to headers in Response.
func SetCookie(name string, value string, c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(name, value, maxAge, "/", domain, true, true)
}

// UnsetCookie removes cookie.
func UnsetCookie(name string, c *gin.Context) {
	c.SetCookie(name, "", -1, "/", domain, true, true)
}
