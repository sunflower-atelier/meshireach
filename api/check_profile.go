package api

import (
	"meshireach/db_model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// CheckProfile はUserIDに対応するデータがあるか確認するメソッド、結果はJSON形式で返す
func CheckProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user db_model.User
		uid, _ := c.Get("UserID")
		db.Where("user_id = ? ", uid).First(&user)
		if user.UserID == uid {
			c.JSON(http.StatusOK, gin.H{
				"status":   "Exist",
				"UserID":   user.UserID,
				"SearchID": user.SearchID,
				"Name":     user.Name,
				"Message":  user.Message,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status":   "NotExist",
				"UserID":   user.UserID,
				"SearchID": user.SearchID,
				"Name":     user.Name,
				"Message":  user.Message,
			})
		}
	}
}
