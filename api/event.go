package api

import (
	"log"
	"meshireach/db/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type result struct {
	Owner    uint      `gorm:"column:event_owner"`
	Title    string    `gorm:"column:event_title"`
	Deadline time.Time `gorm:"column:event_deadline"`
	EventID  uint      `gorm:"column:id"`
}

// GetAllFriendEvents 友達の飯募集を全取得
func GetAllFriendEvents(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := model.User{FirebaseID: c.MustGet("FirebaseID").(string)}

		// UserIDをセット
		db.First(&user)

		var results []result
		db.Table("friendships").Where("user_id = ?", user.ID).Select("events.event_owner, events.event_title, events.event_deadline, events.id").Joins("right join events on events.event_owner = friendships.friend_id").Scan(&results)

		log.Printf("results=%v", results)

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"events": results,
		})
	}
}
