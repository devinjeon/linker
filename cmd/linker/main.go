package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/devinjeon/linker/internal/handlers/auth"
	"github.com/devinjeon/linker/internal/handlers/links"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

const sessionName = "session"

func mustGetEnv(key string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	panic(fmt.Errorf("The environment vairable not set: %s", key))
}

// For development
var isDev = os.Getenv("IS_DEV") == "true"
var devPort = os.Getenv("DEV_PORT")
var endpoint = os.Getenv("DYNAMODB_ENDPOINT")

// App main URL
var linkerURL = mustGetEnv("LINKER_URL")

// Used to validate cookies with HMAC
var sessionSecretKey = mustGetEnv("SESSION_SECRET_KEY")

// Used to encrypt cookies with AES-256. Should be 32 bytes
var sessionEncryptionKey = mustGetEnv("SESSION_ENCRYPTION_KEY")

// OAuth2
var oauth2ClientID = mustGetEnv("OAUTH2_CLIENT_ID")
var oauth2ClientSecret = mustGetEnv("OAUTH2_CLIENT_SECRET")

// DynamoDB Table name to store links
var tableName = mustGetEnv("DYNAMODB_TABLE_NAME")

var r *gin.Engine
var ginLambda *ginadapter.GinLambda

func init() {
	r = gin.Default()

	parsedURL, err := url.Parse(linkerURL)
	if err != nil {
		panic(err)
	}

	// Use middleware for session
	store := cookie.NewStore([]byte(sessionSecretKey), []byte(sessionEncryptionKey))
	options := sessions.Options{
		Path:     "/",
		Domain:   parsedURL.Hostname(),
		MaxAge:   3600 * 24 * 90,
		Secure:   !isDev,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	store.Options(options)
	r.Use(sessions.Sessions(sessionName, store))

	// CORS
	config := cors.DefaultConfig()
	config.AllowCredentials = true
	config.AllowOrigins = []string{linkerURL}
	config.MaxAge = time.Hour * 24 * 90
	r.Use(cors.New(config))

	// Auth
	// Create OAuth2 config
	oauth2Conf := oauth2.Config{
		ClientID:     oauth2ClientID,
		ClientSecret: oauth2ClientSecret,
		Scopes:       []string{},
		Endpoint:     github.Endpoint,
	}

	// Route
	// 1. Auth Handlers
	authHandlers := auth.New(linkerURL, oauth2Conf)
	authGroup := r.Group("/auth")
	{
		authGroup.GET("/signin", authHandlers.SignIn)
		authGroup.GET("/exchange", authHandlers.Exchange)
		authGroup.GET("/signout", authHandlers.SignOut)
	}

	// 2. Links Handlers
	linksHandlers := links.New(tableName, endpoint)

	r.GET("/links/:id", linksHandlers.Redirect)

	linksGroup := r.Group("/links")
	linksGroup.Use(authHandlers.RequireAuth())
	{
		linksGroup.POST("/:id", linksHandlers.Upsert)
		linksGroup.PUT("/:id", linksHandlers.Upsert)
		linksGroup.DELETE("/:id", linksHandlers.Delete)
	}

	ginLambda = ginadapter.New(r)
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	if isDev {
		if devPort == "" {
			devPort = "8081"
		}
		r.Run(fmt.Sprintf(":%s", devPort))
	} else {
		lambda.Start(handler)
	}
}
