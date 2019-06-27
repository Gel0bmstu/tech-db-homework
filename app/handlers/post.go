package handlers

import (
	"encoding/json"
	"errors"
	"fmt"

	"../models"
	"github.com/gin-gonic/gin"
)

// GET

func PostIdDetailsGet(c *gin.Context) {
	var p models.Post

	id := c.Param("id")

	urlParams := c.Request.URL.Query()

	related := urlParams.Get("related")
	fmt.Println("log PD: ", c.FullPath())

	data, err := p.GetPostDetails(id, related)

	switch err {
	case nil:
		sendData(c, 200, data)
	case models.PostNotFoundError:
		sendError(c, 404, err.Error())
	case models.UserNotFoundError:
		sendError(c, 404, err.Error())
	case models.ThreadNotFoundError:
		sendError(c, 404, err.Error())
	case models.ForumNotFoundError:
		sendError(c, 404, err.Error())
	default:
		panic(err)
	}
}

// POST

func PostIdDetailsPost(c *gin.Context) {
	var p models.Post

	id := c.Param("id")

	err := json.NewDecoder(c.Request.Body).Decode(&p)

	if err != nil {
		panic(errors.New("Error: Failed to decode request\n"))
	}

	err = p.PostUpdate(id)

	fmt.Println("edit: ", p.IsEdited)

	switch err {
	case nil:
		sendData(c, 200, p)
	case models.PostNotFoundError:
		sendError(c, 404, err.Error())
	default:
		panic(err)
	}
}
