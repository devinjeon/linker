package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func (c *DB) marshalItem(item interface{}) (map[string]*dynamodb.AttributeValue, error) {
	marshaled, err := dynamodbattribute.MarshalMap(item)
	return marshaled, err
}

func (c *DB) unmarshalItem(marshaledItem map[string]*dynamodb.AttributeValue, unmarshaled interface{}) error {
	err := dynamodbattribute.UnmarshalMap(marshaledItem, &unmarshaled)
	return err
}

// GetItem finds and returns item by key.
func (c *DB) GetItem(key map[string]*dynamodb.AttributeValue, item interface{}) error {
	result, err := c.client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(c.tableName),
		Key:       key,
	})

	if err != nil {
		return ErrDBOperation
	}

	if result.Item == nil {
		return ErrNotFoundItem
	}

	err = c.unmarshalItem(result.Item, item)
	if err != nil {
		return ErrUnmarshalling
	}

	return nil
}

// PutItem create a new item or replace item with new one.
func (c *DB) PutItem(item interface{}) error {
	marshaledItem, err := c.marshalItem(item)
	if err != nil {
		return ErrMarshalling
	}

	_, err = c.client.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(c.tableName),
		Item:      marshaledItem,
	})

	if err != nil {
		return ErrDBOperation
	}

	return nil
}

// DeleteItem delete a item by key.
func (c *DB) DeleteItem(key map[string]*dynamodb.AttributeValue) error {
	err := c.GetItem(key, nil)
	if err != nil {
		return err
	}

	_, err = c.client.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(c.tableName),
		Key:       key,
	})

	if err != nil {
		return ErrDBOperation
	}

	return nil
}
