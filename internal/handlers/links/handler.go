package links

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	m "github.com/devinjeon/linker/internal/middleware"
	db "github.com/devinjeon/linker/internal/utils/dynamodb"

	"github.com/gin-gonic/gin"
)

type newLink struct {
	URL string `json:"url"`
	TTL int    `json:"ttl"`
}

// Redirect is a handler returning 301 redirection to URL named by ID.
func Redirect(c *gin.Context) {
	id := c.Param("id")

	url, err := getURL(id)
	switch err {
	case db.ErrDBOperation:
		c.Status(500)
	case db.ErrNotFoundItem:
		c.Status(404)
	case db.ErrUnmarshalling:
		c.Status(500)
	case nil:
		c.Status(301)
		c.Header("Location", url)
	}
}

// Upsert create or overwrites link.
var Upsert = m.RequireSession(upsert)

func upsert(c *gin.Context) {
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

	sess, ok := m.GetSession(c)
	if !ok {
		c.Status(401)
		return
	}

	user, _ := sess.UserEmail()
	isOwner, err := verifyLinkOwner(id, user)
	if err != nil && err != db.ErrNotFoundItem {
		fmt.Println(err.Error())
		c.Status(500)
		return
	}
	if err == nil && !isOwner {
		c.Status(401)
		return
	}

	err = putURL(id, data.URL, user, data.TTL)
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
var Delete = m.RequireSession(delete)

func delete(c *gin.Context) {
	id := c.Param("id")

	sess, ok := m.GetSession(c)
	if !ok {
		c.Status(401)
		return
	}

	user, _ := sess.UserEmail()
	isOwner, err := verifyLinkOwner(id, user)
	if err != nil {
		fmt.Println(err.Error())
		c.Status(500)
		return
	}
	if !isOwner {
		c.Status(401)
		return
	}

	err = deleteURL(id)
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
