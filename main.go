package main

import (
	"meshireach/api"
	"meshireach/db/model"
	"meshireach/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func initLocale() {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		loc = time.FixedZone("Asia/Tokyo", 9*60*60)
	}
	time.Local = loc
}

func initDB() *gorm.DB {
	db, err := gorm.Open("mysql", "root:@tcp(db:3306)/meshireach?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Event{})

	return db
}

func initRoute(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.AutoOptions())

	r.GET("/ping", api.Ping())

	authedGroup := r.Group("/")
	authedGroup.Use(middleware.FirebaseAuth())
	{
		authedGroup.GET("/private", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "認証成功！")
		})

		// profile
		authedGroup.GET("/profiles", api.CheckProfile(db))
		authedGroup.POST("/profiles", api.CreateProfile(db))
		authedGroup.PUT("/profiles", api.EditProfile(db))

		// friends
		authedGroup.POST("/friends", api.RegisterFriends(db))
		authedGroup.GET("/friends", api.GetAllFriends(db))

		// event
		authedGroup.POST("/events", api.RegisterEvent(db))
	}

	return r
}

func main() {
	initLocale()

	db := initDB()
	defer db.Close()

	r := initRoute(db)
	r.Run(":3000")
}
