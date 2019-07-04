package main

import (
	"meshireach/api"
	"meshireach/db/model"
	"meshireach/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func initDB() *gorm.DB {
	db, err := gorm.Open("mysql", "root:@tcp(db:3306)/meshireach?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Friend_Relation{})
	db.AutoMigrate(&model.Event{})
	db.AutoMigrate(&model.Participants{})

	return db
}

func initRoute(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.AutoOptions())

	r.GET("/ping", api.Ping())

	authedGroup := r.Group("/")
	authedGroup.Use(middleware.FirebaseAuth())
	{
		r.GET("/private", middleware.FirebaseAuth(), func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "認証成功！")
		})

		// profile
		r.POST("/profiles", api.CreateProfile(db))
		r.PUT("/profiles", api.EditProfile(db))
	}

	return r
}

func main() {
	db := initDB()
	defer db.Close()

	r := initRoute(db)
	r.Run(":3000")
}
