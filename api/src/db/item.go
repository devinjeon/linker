package db

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func (c *DB) marshalLinkItem(link Link) (map[string]*dynamodb.AttributeValue, error) {
	marshaledLink, err := dynamodbattribute.MarshalMap(link)
	if err != nil {
		fmt.Printf("Failed to marshal Link item, %v", err)
	}
	return marshaledLink, err
}

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
		TableName: aws.String(c.tableName),
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

// PutURL create a new item or replace item if the id already exists
func (c *DB) PutURL(id string, url string) error {
	newLink := Link{
		ID:  id,
		URL: url,
	}

	marshaledLink, err := c.marshalLinkItem(newLink)
	if err != nil {
		return ErrMarshalling
	}

	_, err = c.client.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(c.tableName),
		Item:      marshaledLink,
	})

	if err != nil {
		fmt.Println(err.Error())
		return ErrDBOperation
	}

	return nil
}

// DeleteURL delete a item from id
func (c *DB) DeleteURL(id string) error {
	_, err := c.GetURL(id)
	if err != nil {
		return err
	}

	_, err = c.client.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(c.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(id),
			},
		},
	})

	if err != nil {
		fmt.Println(err.Error())
		return ErrDBOperation
	}

	return nil
}
