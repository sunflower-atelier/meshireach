package api

import (
	"meshireach/db/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type reqRegister {
	Title    string    `json:"title"`     // タイトル
	Deadline time.Time `json:"deadline"`  // 開始時刻(RFC3339)
}

func validDeadline(t time.Time) bool {
	lower := time.Now()
	upper := time.Truncate(24 * time.Hour).AddDate(0, 0, 3)

	if !t.Before(lower) {
		return false
	}

	if !t.After(upper) {
		return false
	}

	return true
}

func RegisterEvent(db *gorm.DB) gin.HandlerFunc {
	return func (c *gin.Context) {
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
		var owner User
		firebaseid := c.MustGet("FirebaseID").(string)
		db.Where(&User{FirebaseID: firebaseid}).First(&owner)
		db.Create(&Event{Owner: owner.ID, Title: req.Title, Deadline: req.Deadline})
	}
}