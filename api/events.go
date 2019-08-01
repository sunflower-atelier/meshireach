package api

import (
	"meshireach/db/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// RegisterEvents イベントへの参加
func RegisterEvents(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		firebaseid := c.MustGet("FirebaseID").(string)
		eventid := c.MustGet("ID").(string)

		// 登録ユーザを取得(必ず取得できるはず)
		user := model.User{}
		db.Where(&model.User{FirebaseID: firebaseid}).First(&user)

		// イベントを取得(取得できるかわからない)
		event := model.Event{}
		if db.Where(&model.Event{ID: eventid}).First(&event).RecordNotFound() {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"error":  "This event has not be registered.",
			})
			return
		}

		//　既に登録されている場合登録できない
		usercopy := user
		if db.Model(&event).Association("Users").Find(&usercopy).Count() != 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"error":  "Cannot register a user for some events already.",
			})
			return
		}

		// ユーザをイベントに登録
		db.Model(&event).Association("Users").Append(&user)

		// 必要な情報を返す
		c.JSON(http.StatusCreated, gin.H{
			"status":   "success",
			"searchID": user.SearchID,
			"name":     user.Name,
			"message":  user.Message,
		})

	}
}

type friendInfo struct {
	Name     string `json:"name"`
	SearchID string `json:"searchID"`
	Message  string `json:"message"`
}

// GetAllAttendees イベント参加者全取得
func GetAllAttendees(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var event model.Event
		var attendees []model.User
		db.Where(&model.Event{id: c.MustGet("FirebaseID").(string)}).First(&event)
		db.Model(&event).Related(&attendees, "Users")

		// 必要な情報のみをコピー
		result := make([]friendInfo, len(attendees))
		for idx, a := range attendees {
			result[idx].Name = a.Name
			result[idx].SearchID = a.SearchID
			result[idx].Message = a.Message
		}

		// 結果を返す
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"friends": result,
		})
	}
}
