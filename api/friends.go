package api

import (
	"meshireach/db/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type reqRegister struct {
	SearchID string `json:"searchID"`
}

// RegisterFriends 友達登録する
func RegisterFriends(db *gorm.DB) gin.HandlerFunc {
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
		tocopy := to
		if db.Model(&from).Association("Friends").Find(&tocopy).Count() != 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"error":  "Cannot make a friend with some friends already.",
			})
			return
		}


		db.Model(&from).Association("Friends").Append(&to)
		db.Model(&to).Association("Friends").Append(&from)

		c.JSON(http.StatusCreated, gin.H{
			"status":   "success",
			"searchID": to.SearchID,
			"name":     to.Name,
			"message":  to.Message,
		})

	}
}