package links

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	db "github.com/devinjeon/linker/internal/utils/dynamodb"
	"os"
)

var tableName = os.Getenv("DYNAMODB_LINK_TABLE_NAME")
var c = db.NewDB(tableName)

// Link struct is a model for items from DynamoDB table.
type Link struct {
	ID    string `dynamodbav:"id"`
	URL   string `dynamodbav:"url"`
	Owner string `dynamodbav:"owner"`
}

// LinkWithTTL struct is a model for items having TTL from DynamoDB table.
type LinkWithTTL struct {
	ID    string `dynamodbav:"id"`
	URL   string `dynamodbav:"url"`
	Owner string `dynamodbav:"owner"`
	TTL   int    `dynamodbav:"ttl"`
}

// GetURL finds and returns URL from ID in DynamoDB.
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

// PutURL create a new item or replace item if the ID already exists.
// if TTL is 0, it means not set TTL.
func putURL(id string, url string, user string, TTL int) error {
	var newLink interface{}
	if TTL == 0 {
		newLink = Link{
			ID:    id,
			URL:   url,
			Owner: user,
		}
	} else {
		newLink = LinkWithTTL{
			ID:    id,
			URL:   url,
			Owner: user,
			TTL:   TTL,
		}
	}
	err := c.PutItem(newLink)
	return err
}

func verifyLinkOwner(id string, user string) (bool, error) {
	key := map[string]*dynamodb.AttributeValue{
		"id": {
			S: aws.String(id),
		},
	}
	link := LinkWithTTL{}
	err := c.GetItem(key, &link)
	if err != nil {
		return false, err
	}
	return user == link.Owner, nil
}

// DeleteURL deletes a item by ID.
func deleteURL(id string) error {
	key := map[string]*dynamodb.AttributeValue{
		"id": {
			S: aws.String(id),
		},
	}

	err := c.DeleteItem(key)

	if err != nil {
		return db.ErrDBOperation
	}

	return nil
}
