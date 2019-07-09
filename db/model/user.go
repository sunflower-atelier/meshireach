package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	FirebaseID string
	SearchID   string `json:"searchID"`
	Name       string `json:"name"`
	Message    string `json:"message"`
}
