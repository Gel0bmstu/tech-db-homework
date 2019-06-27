package handlers

import (
	"encoding/json"
	"errors"

	"../models"
	"github.com/gin-gonic/gin"
)

// GET

func ForumSlugDetails(c *gin.Context) {
	var f models.Forum

	slug := c.Param("slug")

	f.Slug = slug

	err := f.GetForumDetails()

	switch err {
	case nil:
		sendData(c, 200, f)
	case models.ForumNotFoundError:
		sendError(c, 404, err.Error())
	default:
		panic(err)
	}
}

func ForumSlugUsers(c *gin.Context) {
	var (
		u  models.Users
		uv models.UrlVars
	)

	slug := c.Param("slug")

	urlParams := c.Request.URL.Query()

	uv.Limit = urlParams.Get("limit")
	uv.Since = urlParams.Get("since")
	uv.Desc = urlParams.Get("desc")

	if uv.Desc == "true" {
		uv.Desc = "DESC"
	}

	err := u.GetUsers(uv, slug)

	if u == nil {
		votar := models.Users{}
		u = votar
	}

	switch err {
	case nil:
		sendData(c, 200, u)
	case models.ForumNotFoundError:
		sendError(c, 404, err.Error())
	default:
		panic(err)
	}
}

func ForumSlugThreads(c *gin.Context) {
	var (
		f  models.Forum
		uv models.UrlVars
	)

	slug := c.Param("slug")

	f.Slug = slug

	urlParams := c.Request.URL.Query()

	uv.Limit = urlParams.Get("limit")
	uv.Since = urlParams.Get("since")
	uv.Desc = urlParams.Get("desc")

	if uv.Desc == "true" {
		uv.Desc = "DESC"
	}

	threads, err := f.GetThreads(uv)

	if threads == nil {
		votar := models.Threads{}
		threads = votar
	}

	switch err {
	case nil:
		sendData(c, 200, threads)
	case models.ForumNotFoundError:
		sendError(c, 404, err.Error())
	default:
		panic(err)
	}
}

// POST

func ForumCreate(c *gin.Context) {
	var f models.Forum

	err := json.NewDecoder(c.Request.Body).Decode(&f)

	if err != nil {
		panic(errors.New("Error: Failed to decode request\n"))
	}

	err = f.CreateForum()

	switch err {
	case nil:
		sendData(c, 201, f)
	case models.ForumConflictError:
		sendData(c, 409, f)
	case models.ForumUserNotFoundError:
		sendError(c, 404, err.Error())
	default:
		panic(err)
	}
}

func ForumSlugCreate(c *gin.Context) {
	var t models.Thread

	err := json.NewDecoder(c.Request.Body).Decode(&t)

	if err != nil {
		panic(errors.New("Error: Failed to decode request\n"))
	}

	forumSlug := c.Param("slug")

	t.Forum = forumSlug

	err = t.CreateThread()

	switch err {
	case nil:
		sendData(c, 201, t)
	case models.ThreadAlreadyExistError:
		sendData(c, 409, t)
	case models.ThreadAuthorNotFoundError:
		sendError(c, 404, err.Error())
	case models.ThreadForumNotFoundError:
		sendError(c, 404, err.Error())
	default:
		panic(err)
	}
}
