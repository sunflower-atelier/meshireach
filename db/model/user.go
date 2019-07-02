package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	UserID   string
	SearchID string
	Name     string
	Message  string
}
