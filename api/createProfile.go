package api

import (
	"fmt"
	"meshireach/db/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// CreateProfile is called when create user information
func CreateProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := model.User{}
		c.BindJSON(&user)

		if uniqueSearchID(db, user.SearchID) {
			firebaseID := c.MustGet("FirebaseID")
			user.FirebaseID = firebaseID.(string)
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
