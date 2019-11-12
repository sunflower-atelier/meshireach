package main

import (
	"context"
	"fmt"
	"meshireach/api"
	"meshireach/db/model"
	"meshireach/middleware"
	"net/http"
	"os"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"google.golang.org/api/option"
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
	db.AutoMigrate(&model.Device{})

	return db
}

func initFirebaseApp() *firebase.App {
	opt := option.WithCredentialsFile("./key/otameshi-firebase-adminsdk.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		fmt.Printf("error initializing app: %v", err)
		os.Exit(1)
	}
	return app
}

func initRoute(db *gorm.DB, fapp *firebase.App) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.AutoOptions())

	// 生存確認用
	r.GET("/ping", api.Ping())

	// 通知テスト用
	r.GET("/messaging", testMessaging(fapp))

	authedGroup := r.Group("/")
	authedGroup.Use(middleware.FirebaseAuth(fapp))
	{
		authedGroup.GET("/private", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "認証成功！")
		})

		// profile
		authedGroup.GET("/profiles", api.CheckProfile(db))
		authedGroup.POST("/profiles", api.CreateProfile(db))
		authedGroup.PUT("/profiles", api.EditProfile(db))

		// friends
		authedGroup.POST("/friends", api.RegisterFriends(fapp, db))
		authedGroup.GET("/friends", api.GetAllFriends(db))

		// events
		authedGroup.POST("/events/:id/join", api.JoinEvents(db))
		authedGroup.GET("/events/:id/participants", api.GetAllParticipants(db))
		authedGroup.POST("/events", api.RegisterEvent(fapp, db))
		authedGroup.GET("/events", api.GetMyEvents(db))

		// subscriptions
		authedGroup.GET("/events-subscriptions", api.GetAllFriendEvents(db))

		// device
		authedGroup.POST("/device/token", api.RegisterDeviceToken(db))
    
		// joining-list
		authedGroup.GET("/events-joining-list", api.GetAllJoinEvents(db))

	}

	return r
}

func testMessaging(app *firebase.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtain a messaging.Client from the App.
		netctx := context.Background()
		client, err := app.Messaging(netctx)
		if err != nil {
			s := fmt.Sprintf("error getting Messaging client: %v\n", err)
			fmt.Print(s)
			c.String(http.StatusOK, s)
			c.Abort()
			return
		}

		// This registration token comes from the client FCM SDKs.
		registrationToken := c.Query("token")
		fmt.Printf("[testMessaging] get token: %v", registrationToken)

		// See documentation on defining a message payload.
		message := &messaging.Message{
			Notification: &messaging.Notification{
				Title: "TESTテストtest",
				Body:  "HOGEほげhoge",
			},
			Token: registrationToken,
		}

		// Send a message to the device corresponding to the provided
		// registration token.
		response, err := client.Send(netctx, message)
		if err != nil {
			s := fmt.Sprintf("error response Messaging client: %v\n", err)
			fmt.Print(s)
			c.String(http.StatusOK, s)
			c.Abort()
			return
		}

		// Response is a message ID string.
		c.String(http.StatusOK, fmt.Sprintf("Successfully sent message:%v\n", response))
	}
}

func main() {
	initLocale()

	db := initDB()
	defer db.Close()

	app := initFirebaseApp()

	r := initRoute(db, app)
	r.Run(":3000")
}
