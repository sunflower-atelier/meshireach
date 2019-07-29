package api

import (
	"meshireach/db/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// GetAllFriendEvents 友達の飯募集を全取得
func GetAllFriendEvents(db *gorm.DB) gin.HandlerFunc {
	type result struct {
		Owner    uint      `gorm:"column:event_owner"`
		Title    string    `gorm:"column:event_title"`
		Deadline time.Time `gorm:"column:event_deadline"`
		EventID  uint      `gorm:"column:id"`
	}

	return func(c *gin.Context) {
		user := model.User{FirebaseID: c.MustGet("FirebaseID").(string)}

		// UserIDをセット
		db.First(&user)

		var results []result
		// join tableのfriendships tableとevents tableをJOINすることで
		// events tableから友達の飯募集だけを抽出
		// + 現在時刻よりあとのもののみを抽出
		db.Table("friendships").Where("user_id = ?", user.ID).
			Select("events.event_owner, events.event_title, events.event_deadline, events.id").
			Joins("right join events on events.event_owner = friendships.friend_id AND events.event_deadline > ?", time.Now()).
			Scan(&results)

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"events": results,
		})
	}
}
