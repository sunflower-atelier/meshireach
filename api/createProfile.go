package api

import (
	"fmt"
	"meshireach/db/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func firstCreate(db *gorm.DB, firebaseID string) bool {
	search := model.User{}
	count := 0

	db.Where("firebase_id = ?", firebaseID).Find(&search).Count(&count)
	if count != 0 {
		return false
	}
	return true
}

// CreateProfile is called when create user information
func CreateProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := model.User{}
		c.BindJSON(&user)

		user.FirebaseID = c.MustGet("FirebaseID").(string)

		// 一度profileをcreateしていたらerror
		if !firstCreate(db, user.FirebaseID) {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "fail",
				"error":  "This user has already created his/her profile.",
			})
			return
		}

		if uniqueSearchID(db, user.SearchID) {
			db.Create(&user)

			c.JSON(http.StatusOK, gin.H{
				"status":   "success",
				"searchID": user.SearchID,
				"name":     user.Name,
				"message":  user.Message,
			})
		} else {
			fmt.Printf("SearchID is not unique\n")

			c.JSON(http.StatusBadRequest, gin.H{
				"status": "fail",
				"error":  "search id is not unique.",
			})
		}

	}
}
