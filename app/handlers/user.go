package handlers

import (
	"encoding/json"
	"errors"

	"../models"

	"github.com/gin-gonic/gin"
)

// GET

func UserNicknameProfileGet(c *gin.Context) {
	var u models.User

	nickname := c.Param("nickname")

	u.Nickname = nickname

	err := u.GetUserInfo()

	switch err {
	case nil:
		sendData(c, 200, u)
	case models.UserNotFoundError:
		sendError(c, 404, err.Error())
	}
}

// POST

func UserNicknameCreate(c *gin.Context) {

	var u models.User

	err := json.NewDecoder(c.Request.Body).Decode(&u)

	if err != nil {
		panic(errors.New("Error: Failed to decode request\n"))
	}

	nickname := c.Param("nickname")

	u.Nickname = nickname

	users, err := u.CreateUser()

	switch err {
	case nil:
		sendData(c, 201, u)
	case models.UserConflictError:
		sendData(c, 409, users)
	default:
		panic(err)
	}

}

func UserNicknameProfilePost(c *gin.Context) {
	var u models.User

	err := json.NewDecoder(c.Request.Body).Decode(&u)

	if err != nil {
		panic(errors.New("Error: Failed to decode request\n"))
	}

	nickname := c.Param("nickname")

	u.Nickname = nickname

	err = u.EditUserInfo()

	switch err {
	case nil:
		sendData(c, 200, u)
	case models.UserConflictError:
		sendError(c, 409, err.Error())
	case models.UserNotFoundError:
		sendError(c, 404, err.Error())
	default:
		panic(err)
	}
}
