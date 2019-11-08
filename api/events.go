package api

import (
	"meshireach/db/model"
	"net/http"
	"strconv"
	"time"

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
		db.Table("users_events").Where("user_id = ? AND event_id = ?", user.ID, event.ID).Count(&userCount)
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
			"status":       "success",
			"owener":       event.Owner,
			"title":        event.Title,
			"deadline":     event.Deadline,
			"participants": result,
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

// DBからeventを取ってくるときに使うstruct
type result struct {
	EventID  uint      `gorm:"column:id"`
	Title    string    `gorm:"column:event_title"`
	OwnerID  uint      `gorm:"column:event_owner"`
	Deadline time.Time `gorm:"column:event_deadline"`
}

// GetAllFriendEvents 友達の飯募集を全取得
func GetAllFriendEvents(db *gorm.DB) gin.HandlerFunc {

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

		sendEventToClient(results, db, c)
	}
}

// GetMyEvents 自分の飯募集を全取得
func GetMyEvents(db *gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
		user := model.User{}
		firebaseID := c.MustGet("FirebaseID").(string)

		// UserIDをセット
		db.Where(&model.User{FirebaseID: firebaseID}).First(&user)

		var results []result
		// events tableから自分の飯募集だけを抽出
		// + 現在時刻よりあとのもののみを抽出
		db.Model(&model.Event{}).Where("event_owner = ? AND event_deadline > ?", user.ID, time.Now()).Scan(&results)

		sendEventToClient(results, db, c)
	}
}

// GetAllJoinEvents 自分が参加している飯募集を全取得
func GetAllJoinEvents(db *gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
		user := model.User{}
		firebaseID := c.MustGet("FirebaseID").(string)

		// UserIDをセット
		db.Where(&model.User{FirebaseID: firebaseID}).First(&user)

		var results []result
		var joinlist []model.Event
		// 自分が参加している飯募集のリストを取得
		db.Model(&user).Related(&joinlist, "Events")
		db.Model(&joinlist).Where("event_deadline > ?", time.Now().Add(-10*time.Minute)).Scan(&results)

		/*
			for i := range joinlist {
				tmp := result{joinlist[i].ID, joinlist[i].Title, joinlist[i].Owner, joinlist[i].Deadline}
				results = append(results, tmp)
			}
		*/

		sendEventToClient(results, db, c)
	}
}

func sendEventToClient(results []result, db *gorm.DB, c *gin.Context) {

	// response用
	type eventPlus struct {
		EventID       uint      `json:"id"`
		Title         string    `json:"title"`
		OwnerSearchID string    `json:"ownerID"`
		OwnerName     string    `json:"owner"`
		Deadline      time.Time `json:"deadline"`
	}

	events := []eventPlus{}
	// 各eventのownerのsearch IDとnameを取得
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
