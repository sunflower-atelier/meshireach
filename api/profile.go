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

// searchIDがdb内で一意であるかを確認
func uniqueSearchID(db *gorm.DB, searchID string) bool {
	search := model.User{}
	count := 0

	db.Where("search_id = ?", searchID).Find(&search).Count(&count)
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
				"status": "failure",
				"error":  "This user has already created his/her profile.",
			})
			return
		}

		if uniqueSearchID(db, user.SearchID) {
			db.Create(&user)

			c.JSON(http.StatusCreated, gin.H{
				"status":   "success",
				"searchID": user.SearchID,
				"name":     user.Name,
				"message":  user.Message,
			})
		} else {
			fmt.Printf("SearchID is not unique\n")

			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"error":  "Search ID is not unique.",
			})
		}

	}
}

// EditProfile is called when edit user information
func EditProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := model.User{}
		c.BindJSON(&user)

		if uniqueSearchID(db, user.SearchID) {
			// 更新前の情報を取得
			firebaseID := c.MustGet("FirebaseID")
			beforeUser := model.User{FirebaseID: firebaseID.(string)}
			db.First(&beforeUser)

			// 更新
			db.Model(&beforeUser).Updates(user)

			c.JSON(http.StatusOK, gin.H{
				"status":   "success",
				"searchID": user.SearchID,
				"name":     user.Name,
				"message":  user.Message,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"error":  "Search ID is not unique.",
			})
		}
	}
}

// CheckProfile はDBにプロフィールがあるか確認するAPI
func CheckProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user model.User
		firebaseid := c.MustGet("FirebaseID")
		if db.Where("firebase_id = ? ", firebaseid).First(&user).RecordNotFound() == false {
			c.JSON(http.StatusOK, gin.H{
				"status":   "Exist",
				"searchID": user.SearchID,
				"name":     user.Name,
				"message":  user.Message,
			})
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"status": "NotExist",
			})
		}
	}
}
