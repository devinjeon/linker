package main

import (
	"context"
	"os"
	"time"

	"github.com/devinjeon/linker/internal/handlers/auth"
	"github.com/devinjeon/linker/internal/handlers/links"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var linkerDomain = os.Getenv("LINKER_DOMAIN")
var ginLambda *ginadapter.GinLambda

func init() {
	r := gin.Default()

	// CORS
	config := cors.DefaultConfig()
	config.AllowCredentials = true
	config.AllowOrigins = []string{linkerDomain}
	config.MaxAge = time.Hour * 24 * 30
	r.Use(cors.New(config))

	// Route
	linksGroup := r.Group("/links")
	{
		linksGroup.GET("/:id", links.Redirect)
		linksGroup.POST("/:id", links.Upsert)
		linksGroup.PUT("/:id", links.Upsert)
		linksGroup.DELETE("/:id", links.Delete)
	}
	authGroup := r.Group("/auth")
	{
		authGroup.GET("/signin", auth.SignIn)
		authGroup.GET("/signout", auth.SignOut)
		authGroup.GET("/exchange", auth.Exchange)
	}

	ginLambda = ginadapter.New(r)
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(handler)
}
