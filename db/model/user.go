package model

import "github.com/jinzhu/gorm"

// User ユーザーのモデル
// Friendsは友達関係(joinテーブル名はfriendsships，相手のidはfriend_idカラムになる)
// Eventsは参加イベントのリスト(joinテーブル名はparticipants)
type User struct {
	gorm.Model
	FirebaseID string   `gorm:"column:firebase_id"`                                                // FirebaseのID(一意)
	SearchID   string   `gorm:"column:search_id" json:"searchID"`                                  // 検索ID(一意)
	Name       string   `gorm:"column:user_name" json:"name"`                                      // スクリーンネーム
	Message    string   `gorm:"column:user_message" json:"message"`                                // 一言メッセージ
	Friends    []*User  `gorm:"many2many:friendships; association_jointable_foreignkey:friend_id"` // 友達リスト
	Events     []*Event `gorm:"many2many:participants"`                                            // 参加しているイベント
	// 作ったイベントも持つべきかも
}
