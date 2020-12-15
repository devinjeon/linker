module auth

go 1.15

require (
	github.com/aws/aws-lambda-go v1.20.0
	linker/utils/http v0.0.0
	linker/utils/oauth2 v0.0.0
)

replace (
	linker/utils/http v0.0.0 => ../../utils/http
	linker/utils/oauth2 v0.0.0 => ../../utils/oauth2
)
