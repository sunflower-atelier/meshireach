package db_model

import "github.com/jinzhu/gorm"

type Friend_Relation struct {
	gorm.Model
	From UserID
	To   UserID
}