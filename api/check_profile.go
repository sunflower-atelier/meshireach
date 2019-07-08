package api

import (
	"meshireach/db/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// CheckProfile はDBにプロフィールがあるか確認するAPI
func CheckProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user model.User
		uid, _ := c.Get("UserID")
		db.Where("user_id = ? ", uid).First(&user)
		if user.ID == uid {
			c.JSON(http.StatusOK, gin.H{
				"status":     "Exist",
				"ID":         user.ID,
				"FirebaseID": user.FirebaseID,
				"SearchID":   user.SearchID,
				"Name":       user.Name,
				"Message":    user.Message,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": "NotExist",
			})
		}
	}
}
