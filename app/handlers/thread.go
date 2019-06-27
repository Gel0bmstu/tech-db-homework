package handlers

import (
	"encoding/json"
	"errors"

	"../models"

	"github.com/gin-gonic/gin"
)

// GET

func ThreadsSlugOrIdDetailsGet(c *gin.Context) {
	var t models.Thread

	slugOrID := c.Param("slugOrID")

	err := t.GetThread(slugOrID)

	switch err {
	case nil:
		sendData(c, 200, t)
	case models.ThreadNotFoundError:
		sendError(c, 404, err.Error())
	default:
		panic(err)
	}
}

func ThreadsSlugOrIdPosts(c *gin.Context) {
	var (
		p  models.Posts
		uv models.UrlVars
		// posts *models.Posts
	)

	slugOrID := c.Param("slugOrID")

	urlParams := c.Request.URL.Query()

	uv.Limit = urlParams.Get("limit")
	uv.Since = urlParams.Get("since")
	uv.Desc = urlParams.Get("desc")
	uv.Sort = urlParams.Get("sort")

	if uv.Desc == "true" {
		uv.Desc = "DESC"
	}

	posts, err := p.GetPosts(slugOrID, uv)

	// res, _ := json.Marshal(posts)

	switch err {
	case nil:
		if *posts == nil {
			votar := []models.Posts{}
			sendData(c, 200, votar)
		} else {
			sendData(c, 200, posts)
		}
	case models.ThreadNotFoundError:
		sendError(c, 404, err.Error())
	default:
		panic(err)
	}
}

// POST

func ThreadsSlugOrIdCrete(c *gin.Context) {
	var posts, p models.Posts

	err := json.NewDecoder(c.Request.Body).Decode(&p)

	if err != nil {
		panic(errors.New("Error: Failed to decode request\n"))
	}

	slugOrID := c.Param("slugOrID")

	err = p.CreatePost(slugOrID)

	if posts == nil {
		votar := models.Posts{}
		posts = votar
	}

	switch err {
	case nil:
		sendData(c, 201, p)
	case models.ParentPostNotFoundInThread:
		sendError(c, 409, err.Error())
	case models.ThreadNotFoundError:
		sendError(c, 404, err.Error())
	case models.UserNotFoundError:
		sendError(c, 404, err.Error())
	default:
		panic(err)
	}
}

func ThreadsSlugOrIdDetailsPost(c *gin.Context) {
	var t models.Thread

	err := json.NewDecoder(c.Request.Body).Decode(&t)

	if err != nil {
		panic(errors.New("Error: Failed to decode request\n"))
	}

	slugOrID := c.Param("slugOrID")

	err = t.ThreadUpdate(slugOrID)

	switch err {
	case nil:
		sendData(c, 200, t)
	case models.ThreadNotFoundError:
		sendError(c, 404, err.Error())
	default:
		panic(err)
	}
}

func ThreadsSlugOrIdVote(c *gin.Context) {
	var (
		v models.Vote
		t *models.Thread
	)

	err := json.NewDecoder(c.Request.Body).Decode(&v)

	if err != nil {
		panic(errors.New("Error: Failed to decode request\n"))
	}

	slugOrID := c.Param("slugOrID")

	t, err = v.ChangeVote(slugOrID)

	switch err {
	case nil:
		sendData(c, 200, t)
	case models.ThreadNotFoundError:
		sendError(c, 404, err.Error())
	case models.UserNotFoundError:
		sendError(c, 404, err.Error())
	default:
		panic(err)
	}
}
