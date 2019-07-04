package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Event struct {
	gorm.Model
	EventID  int
	Owner    int
	Title    string
	Deadline time.Time
}
