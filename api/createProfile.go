package api

import (
	"meshireach/db/model"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// CreateProfile is called when create user information
func CreateProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user model.User
		c.BindJSON(&user)
		if firebaseID, exists := c.Get("FirebaseID"); exists {
			user.FirebaseID = firebaseID.(string)
			db.Create(user)

			c.JSON(200, gin.H{
				"status":   "success",
				"searchID": user.SearchID,
				"name":     user.Name,
				"message":  user.Message,
			})
		} else {
			c.JSON(400, gin.H{
				"status": "fail",
			})
		}

	}
}
