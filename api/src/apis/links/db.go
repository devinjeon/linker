package links

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	db "linker/utils/dynamodb"
	"os"
)

var tableName = os.Getenv("DYNAMODB_LINK_TABLE_NAME")
var c = db.NewDB(tableName)

// Link struct is a model for items from DynamoDB table
type Link struct {
	ID  string `dynamodbav:"id"`
	URL string `dynamodbav:"url"`
}

// LinkWithTTL struct is a model for items having ttl from DynamoDB table
type LinkWithTTL struct {
	ID  string `dynamodbav:"id"`
	URL string `dynamodbav:"url"`
	TTL int    `dynamodbav:"ttl"`
}

// GetURL finds and returns URL from id in DynamoDB
func getURL(id string) (string, error) {
	key := map[string]*dynamodb.AttributeValue{
		"id": {
			S: aws.String(id),
		},
	}
	link := LinkWithTTL{}
	err := c.GetItem(key, &link)
	return link.URL, err
}

// PutURL create a new item or replace item if the id already exists.
// if ttl is 0, it means not set ttl.
func putURL(id string, url string, ttl int) error {
	var newLink interface{}
	if ttl == 0 {
		newLink = Link{
			ID:  id,
			URL: url,
		}
	} else {
		newLink = LinkWithTTL{
			ID:  id,
			URL: url,
			TTL: ttl,
		}
	}
	err := c.PutItem(newLink)
	return err
}

// DeleteURL delete a item from id
func deleteURL(id string) error {
	key := map[string]*dynamodb.AttributeValue{
		"id": {
			S: aws.String(id),
		},
	}

	err := c.DeleteItem(key)

	if err != nil {
		fmt.Println(err.Error())
		return db.ErrDBOperation
	}

	return nil
}
