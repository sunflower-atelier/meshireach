package api

import (
	"meshireach/db/model"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// CheckProfile はDBにプロフィールがあるか確認するAPI
func CheckProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user model.User
		firebaseid := c.MustGet("FirebaseID")
		if db.Where("firebase_id = ? ", firebaseid).First(&user).RecordNotFound() == false {
			c.JSON(200, gin.H{
				"status":     "Exist",
				"ID":         user.ID,
				"FirebaseID": user.FirebaseID,
				"SearchID":   user.SearchID,
				"Name":       user.Name,
				"Message":    user.Message,
			})
		} else {
			c.JSON(404, gin.H{
				"status": "NotExist",
			})
		}
	}
}
