package db_model

import (
	"time"
)

type EventID int

type Event struct {
	EventID  EventID `gorm:"primary_key"`
	Owner    UserID
	Title    string
	Deadline time.Time
}
