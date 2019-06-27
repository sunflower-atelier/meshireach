package main

import (
	"meshireach/db_model"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func editProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user db_model.User
		c.BindJSON(&user)
		c.JSON(200, gin.H{
			"status":  "success",
			"name":    user.Name,
			"message": user.Message,
		})
	}
}

func main() {
	db, err := gorm.Open("mysql", "root:@tcp(db:3306)/meshireach?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&db_model.User{})
	db.AutoMigrate(&db_model.Friend_Relation{})
	db.AutoMigrate(&db_model.Event{})
	db.AutoMigrate(&db_model.Participants{})

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/profile", editProfile(db))
	r.Run(":3000")
}
