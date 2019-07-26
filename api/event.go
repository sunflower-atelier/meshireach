package api

import (
	"log"
	"meshireach/db/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// GetAllFriendEvents 友達の飯募集を全取得
func GetAllFriendEvents(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := model.User{FirebaseID: c.MustGet("FirebaseID").(string)}
		var events []model.Event
		db.Table("friendships").Where("user_id = ?", user.FirebaseID).Joins("right join events on events.event_owner = friendships.friend_id").Scan(&events)

		log.Printf("events[0].Title=%v", events[0].Title)

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"events": events,
		})
	}
}
