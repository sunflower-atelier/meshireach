package db_model

import "github.com/jinzhu/gorm"

type Participants struct {
	gorm.Model
	Event EventID
	User  UserID
}
