package api

import (
	"meshireach/db/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func validDeadline(t time.Time) bool {
	lower := time.Now()
	upper := t.Truncate(24*time.Hour).AddDate(0, 0, 3) // ３日後まで

	if t.After(lower) && t.Before(upper) {
		return true
	} else {
		return false
	}
}

// RegisterEvent イベントの登録
func RegisterEvent(db *gorm.DB) gin.HandlerFunc {
	type reqRegister struct {
		Title    string    `json:"title"`    // タイトル
		Deadline time.Time `json:"deadline"` // 開始時刻(RFC3339)
	}

	return func(c *gin.Context) {
		req := reqRegister{}
		c.BindJSON(&req)

		// 開始時間のバリデーション
		if !validDeadline(req.Deadline) {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"error":  "Invalid Meshi start time.",
			})
			return
		}

		// 登録
		var owner model.User
		firebaseid := c.MustGet("FirebaseID").(string)
		db.Where(&model.User{FirebaseID: firebaseid}).First(&owner)

		event := model.Event{Owner: owner.ID, Title: req.Title, Deadline: req.Deadline}
		db.Create(&event)

		c.JSON(http.StatusCreated, gin.H{
			"eventid": event.ID,
		})
	}
}
