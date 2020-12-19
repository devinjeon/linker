package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))

// DB is a struct to handle DynamoDB
type DB struct {
	tableName string
	client    *dynamodb.DynamoDB
}

// NewDB creates a new client of DynamoDB with the specified table name
func NewDB(tableName string) *DB {
	return &DB{
		tableName: tableName,
		client:    dynamodb.New(sess),
	}
}
