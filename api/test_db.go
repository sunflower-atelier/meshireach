package api

import (
	"fmt"
	"meshireach/db/model"

	"github.com/jinzhu/gorm"
	// Mysql driver
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func initDB() *gorm.DB {
	db, err := gorm.Open("mysql", "root:@tcp(db:3306)/meshireach_test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
		fmt.Println(err)
	}

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Friend_Relation{})
	db.AutoMigrate(&model.Event{})
	db.AutoMigrate(&model.Participants{})

	return db
}
