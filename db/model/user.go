package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	FirebaseID string
	SearchID   string
	Name       string
	Message    string
}
