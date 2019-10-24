package api

import (
	"meshireach/db/model"
	"net/http"
	"context"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"

)


// RegisterDeviceToken デバイストークンとそのオーナーの登録
func RegisterDeviceToken(db *gorm.DB) gin.HandlerFunc {
	type reqRegister struct {
		DeviceToken string `json:"device_token"` // デバイストークン
	}

	return func(c *gin.Context) {
		ownerFirebaseID := c.MustGet("FirebaseID").(string)
		user := model.User{}
		db.Where("firebase_id = ? ", ownerFirebaseID).First(&user)

		req := reqRegister{}
		c.BindJSON(&req)

		db.Create(&model.Device{Owner: user.ID, Token: req.DeviceToken})

		c.JSON(http.StatusCreated, gin.H{
			"status": "success",
		})
	}
}

// SendNotify 与えられたオーナーに紐づいたデバイス全てに通知を送る
func SendNotification(fapp *firebase.App, db *gorm.DB, owners []uint, title string, body string, data *map[string]string) error {
	netctx := context.Background()
	client, err := fapp.Messaging(netctx)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	for _, owner := range owners { // ここもgo routineにできそう
		// get owner's devices
		var devices []model.Device
		db.Where("device_owner = ?", owner).Find(&devices)
		
		for _, dev := range devices {
			wg.Add(1)
			go func(){
				defer wg.Done()
				message := &messaging.Message{
					Data: *data,
					Notification: &messaging.Notification {
						Title: title,
						Body: body,
					},
					Token: dev.Token,
				}
				_, err := client.Send(netctx, message)
				if err != nil {
					//return err
				}
			}()
		}

	}
	wg.Wait()
	return nil
}
