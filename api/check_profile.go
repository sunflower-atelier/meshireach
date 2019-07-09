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
		firebaseid := c.MustGet("FirebaseID")
		if db.Where("firebase_id = ? ", firebaseid).First(&user).RecordNotFound() == false {
			c.JSON(http.StatusOK, gin.H{
				"status":   "Exist",
				"searchID": user.SearchID,
				"name":     user.Name,
				"message":  user.Message,
			})
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"status": "NotExist",
			})
		}
	}
}
