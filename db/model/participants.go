package model

import "github.com/jinzhu/gorm"

type Participants struct {
	gorm.Model
	Event int
	User  int
}
