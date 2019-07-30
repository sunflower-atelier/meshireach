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
		EventID       uint      `gorm:"column:id" json:"id"`
		Title         string    `gorm:"column:event_title" json:"title"`
		OwnerID       uint      `gorm:"column:event_owner"`
		OwnerSearchID string    `gorm:column:search_id json:ownerID`
		OwnerName     string    `gorm:"column:user_name" json:owner`
		Deadline      time.Time `gorm:"column:event_deadline" json:"deadline"`
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

		// 各eventのownerのsearch IDとnameを取得
		// JOINするともうちょい早い気がする
		for i := range results {
			owner := model.User{Model: gorm.Model{ID: results[i].OwnerID}}
			db.First(&owner)
			results[i].OwnerSearchID = owner.SearchID
			results[i].OwnerName = owner.Name
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"events": results,
		})
	}
}
