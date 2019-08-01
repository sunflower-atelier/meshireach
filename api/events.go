package api

import (
	"meshireach/db/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// RegisterEvents イベントへの参加
func RegisterEvents(db *gorm.DB) gin.HandlerFunc {

	type reqRegister struct {
		UserID  uint `json:"userID"`
		EventID uint `json:"eventID"`
	}

	return func(c *gin.Context) {
		req := reqRegister{}
		c.BindJSON(&req)
		userid := req.UserID
		eventid := req.EventID

		// 登録ユーザを取得
		user := model.User{}
		db.Where(&model.User{Model: gorm.Model{ID: userid}}).First(&user)
		if db.Where(&model.User{Model: gorm.Model{ID: userid}}).First(&user).RecordNotFound() {
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

		//　既にイベントに登録されている場合登録できない
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
			"status": "success",
		})

	}
}

// GetAllAttendees イベント参加者全取得
func GetAllAttendees(db *gorm.DB) gin.HandlerFunc {

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
