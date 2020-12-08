module linker

go 1.15

require (
	github.com/aws/aws-lambda-go v1.20.0
	github.com/aws/aws-sdk-go v1.36.3 // indirect
	linker/db v0.0.0
)

replace linker/db v0.0.0 => ./db
