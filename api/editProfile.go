package api

import (
	"meshireach/db_model"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// EditProfile is called when edit user information
func EditProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user db_model.User
		c.BindJSON(&user)
		c.JSON(200, gin.H{
			"status":  "success",
			"name":    user.Name,
			"message": user.Message,
		})
	}
}
