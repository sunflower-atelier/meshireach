package api

import (
	"meshireach/db_model"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// CreateProfile is called when create user information
func CreateProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user db_model.User
		c.BindJSON(&user)
		db.Create(user)

		c.JSON(200, gin.H{
			"status":   "success",
			"searchID": user.SearchID,
			"name":     user.Name,
			"message":  user.Message,
		})
	}
}
