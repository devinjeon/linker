module github.com/devinjeon/linker

go 1.15

require (
	github.com/aws/aws-lambda-go v1.20.0
	github.com/aws/aws-sdk-go v1.36.12
	github.com/awslabs/aws-lambda-go-api-proxy v0.9.0
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-contrib/sessions v0.0.3
	github.com/gin-gonic/gin v1.6.3
	golang.org/x/oauth2 v0.0.0-20180821212333-d2e6202438be
)

replace github.com/devinjeon/linker v0.0.0 => ./
