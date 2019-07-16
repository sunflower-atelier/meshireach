package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Event イベント情報のモデル
// Usersは参加者のリスト(joinテーブル名はparticipants)
type Event struct {
	gorm.Model
	Owner    User      `gorm:"column:event_owner"`     // イベント作成者(joinのためのforeign keyを設定すべきかも)
	Title    string    `gorm:"column:event_title"`     // タイトル
	Deadline time.Time `gorm:"column:event_deadline"`  // 開始時刻
	Users    []*User   `gorm:"many2many:participants"` // 参加者リスト
}
