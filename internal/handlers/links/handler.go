package links

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	auth "github.com/devinjeon/linker/internal/handlers/auth"
	db "github.com/devinjeon/linker/internal/utils/dynamodb"

	"github.com/gin-gonic/gin"
)

type newLink struct {
	URL string `json:"url"`
	TTL int    `json:"ttl"`
}

// Handlers is struct including handler methods.
type Handlers struct {
	table *db.Table
}

// New creates handlers for link resources
func New(dynamoDBTableName string, endpoint string) Handlers {
	hanlders := Handlers{
		table: db.NewTable(dynamoDBTableName, endpoint),
	}

	return hanlders
}

// Redirect is a handler returning 301 redirection to URL named by ID.
func (h *Handlers) Redirect(c *gin.Context) {
	id := c.Param("id")

	url, err := h.getURL(id)
	switch err {
	case db.ErrDBOperation:
		c.Status(500)
	case db.ErrNotFoundItem:
		c.Status(404)
	case db.ErrUnmarshalling:
		c.Status(500)
	case nil:
		c.Redirect(301, url)
	}
}

// Upsert create or overwrites link.
func (h *Handlers) Upsert(c *gin.Context) {
	u := c.MustGet("auth").(auth.Auth)

	id := c.Param("id")
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.Status(400)
	}

	var data newLink
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		fmt.Println(err.Error())
		c.Status(400)
		return
	}

	isOwner, err := h.verifyLinkOwner(id, u.User)
	if err != nil && err != db.ErrNotFoundItem {
		fmt.Println(err.Error())
		c.Status(500)
		return
	}
	if err == nil && !isOwner {
		c.Status(401)
		return
	}

	err = h.putURL(id, data.URL, u.User, data.TTL)
	switch err {
	case db.ErrDBOperation:
		fmt.Println(err.Error())
		c.Status(500)
	case db.ErrMarshalling:
		fmt.Println(err.Error())
		c.Status(500)
	case nil:
		c.Status(204)
	}
}

// Delete removes link.
func (h *Handlers) Delete(c *gin.Context) {
	u := c.MustGet("auth").(auth.Auth)
	id := c.Param("id")

	isOwner, err := h.verifyLinkOwner(id, u.User)
	if err != nil {
		fmt.Println(err.Error())
		c.Status(500)
		return
	}
	if !isOwner {
		c.Status(401)
		return
	}

	err = h.deleteURL(id)
	switch err {
	case db.ErrDBOperation:
		fmt.Println(err.Error())
		c.Status(500)
	case db.ErrNotFoundItem:
		c.Status(404)
	case db.ErrUnmarshalling:
		fmt.Println(err.Error())
		c.Status(500)
	case nil:
		c.Status(204)
	}
}
