module linker

go 1.15

require (
	github.com/aws/aws-lambda-go v1.20.0
	linker/apis/links v0.0.0
	linker/utils/dynamodb v0.0.0
	linker/utils/http v0.0.0
	linker/utils/oauth2 v0.0.0
)

replace (
	linker/apis/links v0.0.0 => ./apis/links
	linker/utils/dynamodb v0.0.0 => ./utils/dynamodb
	linker/utils/http v0.0.0 => ./utils/http
	linker/utils/oauth2 v0.0.0 => ./utils/oauth2
)
