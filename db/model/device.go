package model

import (
	"github.com/jinzhu/gorm"
)

// Device デバイス情報のモデル
type Device struct {
	gorm.Model
	Owner uint   `gorm:"column:device_owner"`
	Token string `gorm:"column:device_token"`
}
