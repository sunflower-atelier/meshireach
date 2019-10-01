package api

import (
	"fmt"
	"meshireach/db/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// RegisterFriends 友達登録する
func RegisterFriends(db *gorm.DB) gin.HandlerFunc {
	type reqRegister struct {
		SearchID string `json:"searchID"`
	}

	return func(c *gin.Context) {
		req := reqRegister{}
		c.BindJSON(&req)
		searchid := req.SearchID
		firebaseid := c.MustGet("FirebaseID").(string)

		// リクエスト送信側を取得(必ず取得できるはず)
		from := model.User{}
		db.Where(&model.User{FirebaseID: firebaseid}).First(&from)

		// 友達になる側を取得(取得できるかわからない)
		to := model.User{}
		if db.Where(&model.User{SearchID: searchid}).First(&to).RecordNotFound() {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"error":  "This user has not be registered.",
			})
			return
		}

		// 自分と友達にはなれない
		if from.FirebaseID == to.FirebaseID {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"error":  "Cannot make a friend with yourself.",
			})
			return
		}

		//　友達と友達にはなれない
		var friendCount = 0 // 1 or 0
		db.Table("friendships").Where("user_id = ? AND friend_id = ?", from.ID, to.ID).Count(&friendCount)
		if friendCount != 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"error":  "Cannot make a friend with some friends already.",
			})
			return
		}

		// 友達関係を登録
		db.Model(&from).Association("Friends").Append(&to)
		db.Model(&to).Association("Friends").Append(&from)

		// 必要な情報を返す
		c.JSON(http.StatusCreated, gin.H{
			"status":   "success",
			"searchID": to.SearchID,
			"name":     to.Name,
			"message":  to.Message,
		})

	}
}

// GetAllFriends 友達情報の全取得
func GetAllFriends(db *gorm.DB) gin.HandlerFunc {
	type friendInfo struct {
		Name     string `json:"name"`
		SearchID string `json:"searchID"`
		Message  string `json:"message"`
	}

	return func(c *gin.Context) {
		// firebaseidからユーザーを取得
		// 取得したユーザーに関連した友達を全取得
		var user model.User
		var friends []model.User
		db.Where(&model.User{FirebaseID: c.MustGet("FirebaseID").(string)}).First(&user)
		db.Model(&user).Related(&friends, "Friends")

		// 必要な情報のみをコピー
		result := make([]friendInfo, len(friends))
		for idx, f := range friends {
			result[idx].Name = f.Name
			result[idx].SearchID = f.SearchID
			result[idx].Message = f.Message
		}

		// 結果を返す
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"friends": result,
		})
	}
}
