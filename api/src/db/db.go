package db

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

// GetDB creates a new client of DynamoDB with the specified table name
func GetDB(tableName string) *DB {
	return &DB{
		tableName: tableName,
		client:    dynamodb.New(sess),
	}
}
