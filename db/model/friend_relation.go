package model

import "github.com/jinzhu/gorm"

type Friend_Relation struct {
	gorm.Model
	From uint
	To   uint
}
