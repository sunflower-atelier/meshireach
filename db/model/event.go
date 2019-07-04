package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Event struct {
	gorm.Model
	Owner    uint
	Title    string
	Deadline time.Time
}
