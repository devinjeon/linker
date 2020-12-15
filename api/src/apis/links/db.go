package links

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	db "linker/utils/dynamodb"
	"os"
)

var tableName = os.Getenv("DYNAMODB_LINKS_TABLE_NAME")
var c = db.NewDB(tableName)

// Link struct is a model for items from DynamoDB table
type Link struct {
	ID  string
	URL string
}

// GetURL finds and returns URL from id in DynamoDB
func getURL(id string) (string, error) {
	key := map[string]*dynamodb.AttributeValue{
		"ID": {
			S: aws.String(id),
		},
	}
	link := Link{}
	err := c.GetItem(key, &link)
	return link.URL, err
}

// PutURL create a new item or replace item if the id already exists
func putURL(id string, url string) error {
	newLink := Link{
		ID:  id,
		URL: url,
	}
	err := c.PutItem(newLink)
	return err
}

// DeleteURL delete a item from id
func deleteURL(id string) error {
	key := map[string]*dynamodb.AttributeValue{
		"ID": {
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
