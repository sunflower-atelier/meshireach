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

		// 更新前の情報を取得
		// FIXME: UserIDはheaderのtokenから取得する
		var beforeUser = db_model.User{UserID: user.UserID}
		db.First(&beforeUser)

		// 更新
		db.Model(&beforeUser).Updates(user)

		c.JSON(200, gin.H{
			"status":   "success",
			"searchID": user.SearchID,
			"name":     user.Name,
			"message":  user.Message,
		})
	}
}
