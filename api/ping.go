package api

import (
	"github.com/gin-gonic/gin"
)

// Ping ピン
func Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	}
}
