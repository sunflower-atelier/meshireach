package api

import (
	"meshireach/db/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// JoinEvents イベントへの参加
func JoinEvents(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idstring := c.Param("id")
		idInt, _ := strconv.Atoi(idstring)
		eventid := uint(idInt)

		// 登録ユーザ(参加者)を取得
		userFirebaseid := c.MustGet("FirebaseID")
		user := model.User{}
		if db.Table("users").Where("firebase_id = ?", userFirebaseid).First(&user).RecordNotFound() {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"error":  "This user has not be registered.",
			})
			return
		}

		// イベントを取得
		event := model.Event{}
		if db.Where(model.Event{Model: gorm.Model{ID: eventid}}).First(&event).RecordNotFound() {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"error":  "This event has not be registered.",
			})
			return
		}

		// 参加者とイベントオーナーが友達かどうか確認
		var friendCount = 0
		db.Table("friendships").Where("user_id = ? AND friend_id = ?", user.ID, event.Owner).Count(&friendCount)
		if friendCount == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"error":  "Event's owner is not your friend.",
			})
			return
		}

		//　既にイベントに登録されている場合登録できない
		var userCount = 0
		db.Table("participants").Where("user_id = ? AND event_id = ?", user.ID, event.ID).Count(&userCount)
		if userCount != 0 {
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
			"status": "success",
		})

	}
}

// GetAllParticipants イベント参加者全取得
func GetAllParticipants(db *gorm.DB) gin.HandlerFunc {

	type reqRegister struct {
		EventID uint `json:"eventID"`
	}
	type friendInfo struct {
		Name     string `json:"name"`
		SearchID string `json:"searchID"`
		Message  string `json:"message"`
	}

	return func(c *gin.Context) {
		idstring := c.Param("id")
		idInt, _ := strconv.Atoi(idstring)
		eventid := uint(idInt)

		// イベントを取得
		event := model.Event{}
		if db.Where(model.Event{Model: gorm.Model{ID: eventid}}).First(&event).RecordNotFound() {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"error":  "This event has not be registered.",
			})
			return
		}

		// 参加者リストを取得
		var attendees []model.User
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
			"status":   "success",
			"owener":   event.Owner,
			"title":    event.Title,
			"deadline": event.Deadline,
			"friends":  result,
		})
	}
}
