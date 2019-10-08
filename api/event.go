package api

import (
	"meshireach/db/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// GetMyEvents 自分の飯募集を全取得
func GetMyEvents(db *gorm.DB) gin.HandlerFunc {
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
		user := model.User{}
		firebaseID := c.MustGet("FirebaseID").(string)

		// UserIDをセット
		db.Where(&model.User{FirebaseID: firebaseID}).First(&user)

		var results []result
		// join tableのfriendships tableとevents tableをJOINすることで
		// events tableから友達の飯募集だけを抽出
		// + 現在時刻よりあとのもののみを抽出
		db.Table("friendships").Where("user_id = ?", user.ID).
			Select("events.id, events.event_title, events.event_owner, events.event_deadline").
			Joins("inner join events on events.event_owner = friendships.friend_id AND events.event_deadline > ?", time.Now()).
			Scan(&results)

		events := []eventPlus{}
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
			"events": events,
		})
	}
}

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
		user := model.User{}
		firebaseID := c.MustGet("FirebaseID").(string)

		// UserIDをセット
		db.Where(&model.User{FirebaseID: firebaseID}).First(&user)

		var results []result
		// join tableのfriendships tableとevents tableをJOINすることで
		// events tableから友達の飯募集だけを抽出
		// + 現在時刻よりあとのもののみを抽出
		db.Table("friendships").Where("user_id = ?", user.ID).
			Select("events.id, events.event_title, events.event_owner, events.event_deadline").
			Joins("inner join events on events.event_owner = friendships.friend_id AND events.event_deadline > ?", time.Now()).
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
			"events": events,
		})
	}
}

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
			"eventId": event.ID,
		})
	}
}
