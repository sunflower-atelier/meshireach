package model

import (
	"time"
	"github.com/jinzhu/gorm"
)

type EventID int

type Event struct {
	gorm.Model
	EventID  EventID
	Owner    int
	Title    string
	Deadline time.Time
}
