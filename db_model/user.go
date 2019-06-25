package db_model

type UserID string

type User struct {
	UserID   UserID `gorm:"primary_key"`
	SearchID string
	Name     string
	Message  string
}
