package handlers

import (
	"../models"

	"github.com/gin-gonic/gin"
)

// GET

func ServiceStatus(c *gin.Context) {
	var s models.Service

	err := s.GetService()

	switch err {
	case nil:
		sendData(c, 200, s)
	default:
		panic(err)
	}
}

// POST

func ServiceClear(c *gin.Context) {
	var s models.Service

	err := s.ClearDb()

	switch err {
	case nil:
		sendData(c, 200, s)
	default:
		panic(err)
	}
}
