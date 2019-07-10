package api

import (
	"meshireach/db/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// searchIDがdb内で一意であるかを確認
func uniqueSearchID(db *gorm.DB, searchID string) bool {
	search := model.User{}
	count := 0

	db.Where("search_id = ?", searchID).Find(&search).Count(&count)
	if count != 0 {
		return false
	}
	return true
}

// EditProfile is called when edit user information
func EditProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := model.User{}
		c.BindJSON(&user)

		if uniqueSearchID(db, user.SearchID) {
			// 更新前の情報を取得
			firebaseID := c.MustGet("FirebaseID")
			beforeUser := model.User{FirebaseID: firebaseID.(string)}
			db.First(&beforeUser)

			// 更新
			db.Model(&beforeUser).Updates(user)

			c.JSON(http.StatusOK, gin.H{
				"status":   "success",
				"searchID": user.SearchID,
				"name":     user.Name,
				"message":  user.Message,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "fail",
				"error":  "Search ID is not unique.",
			})
		}
	}
}
