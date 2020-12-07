package db

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func (c *DB) unmarshalLinkItem(marshaledItem map[string]*dynamodb.AttributeValue) (Link, error) {
	item := Link{}
	err := dynamodbattribute.UnmarshalMap(marshaledItem, &item)
	if err != nil {
		fmt.Printf("Failed to unmarshal Record, %v", err)
	}
	return item, err
}

// GetURL finds and returns URL from id in DynamoDB
func (c *DB) GetURL(id string) (string, error) {
	url := ""
	result, err := c.client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(c.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(id),
			},
		},
	})

	if err != nil {
		fmt.Println(err.Error())
		return url, ErrDBOperation
	}

	if result.Item == nil {
		fmt.Println("Could not find '" + id + "'")
		return url, ErrNotFoundItem
	}

	link, err := c.unmarshalLinkItem(result.Item)
	if err != nil {
		fmt.Printf("Failed to unmarshal Record, %v", err)
		return url, ErrUnmarshalling
	}

	url = link.URL
	return url, nil
}
