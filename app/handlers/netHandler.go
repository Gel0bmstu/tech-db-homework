package handlers

import (
	"github.com/gin-gonic/gin"
)

// s - http status, m - error message
func sendError(c *gin.Context, s int, m string) {
	c.JSON(s, gin.H{"error": m})
}

// s - http status, i - data to json encode
func sendData(c *gin.Context, s int, i interface{}) {
	c.JSON(s, i)
}
