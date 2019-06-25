package db_model

import "github.com/jinzhu/gorm"

type UserID string

type User struct {
	gorm.Model
	UserID   UserID
	SearchID string
	Name     string
	Message  string
}
