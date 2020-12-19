package dynamodb

import "errors"

var (
	// ErrNotFoundItem is a error that represents "Cannot find the item"
	ErrNotFoundItem = errors.New("db: cannot find the item")
	// ErrDBOperation is a error that represents "DynamoDB operation failed"
	ErrDBOperation = errors.New("db: DynamoDB operation failed")
	// ErrUnmarshalling is a error that represents "Failed to unmarshal the item""
	ErrUnmarshalling = errors.New("db: failed to unmarshal the item")
	// ErrMarshalling is a error that represents "Failed to unmarshal the item""
	ErrMarshalling = errors.New("db: failed to marshal the item")
)
