package db_model

type UserID string

type User struct {
	UserID   UserID `gorm:"primary_key", json:"userID"`
	SearchID string `json:"searchID"`
	Name     string `json:"name"`
	Message  string `json:"message"`
}
