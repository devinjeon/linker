package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))

// Table is a struct to handle DynamoDB table.
type Table struct {
	tableName string
	client    *dynamodb.DynamoDB
}

// NewTable creates a new client of DynamoDB with the specified table name.
func NewTable(tableName string) *Table {
	return &Table{
		tableName: tableName,
		client:    dynamodb.New(sess),
	}
}
