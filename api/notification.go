package api

import (
	"context"
	"meshireach/db/model"
	"net/http"
	"sync"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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

		count := 0
		db.Table("devices").Where("device_owner = ? AND device_token = ?", user.ID, req.DeviceToken).Count(&count)
		if count != 0 {
			c.JSON(http.StatusOK, gin.H{
				"status": "existed",
			})
			return
		}

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
	for _, owner := range owners {
		wg.Add(1)
		go func(owner uint) {
			defer wg.Done()

			// get owner's devices
			var devices []model.Device
			db.Where("device_owner = ?", owner).Find(&devices)

			for _, dev := range devices {
				wg.Add(1)
				go func(devToken string) {
					defer wg.Done()
					message := &messaging.Message{
						Data: *data,
						Notification: &messaging.Notification{
							Title: title,
							Body:  body,
						},
						Token: dev.Token,
					}
					_, err := client.Send(netctx, message)
					if err != nil {
						//return err
					}
				}(dev.Token)
			}
		}(owner)

	}
	wg.Wait()
	return nil
}
