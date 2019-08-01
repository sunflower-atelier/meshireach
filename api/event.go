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
	// DB用
	type result struct {
		EventID  uint      `gorm:"column:id"`
		Title    string    `gorm:"column:event_title"`
		OwnerID  uint      `gorm:"column:event_owner"`
		Deadline time.Time `gorm:"column:event_deadline"`
	}

	// response用
	type eventPlus struct {
		EventID       uint      `json:"id"`
		Title         string    `json:"title"`
		OwnerSearchID string    `json:"ownerID"`
		OwnerName     string    `json:"owner"`
		Deadline      time.Time `json:"deadline"`
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

		var events []eventPlus
		// 各eventのownerのsearch IDとnameを取得
		// JOINするともうちょい早い気がする
		for i := range results {
			owner := model.User{Model: gorm.Model{ID: results[i].OwnerID}}
			db.First(&owner)

			tmp := eventPlus{results[i].EventID, results[i].Title, owner.SearchID, owner.Name, results[i].Deadline}
			events = append(events, tmp)
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"events": results,
		})
	}
}
