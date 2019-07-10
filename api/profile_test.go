package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func setDummyID() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("FirebaseID", "dummy")
		c.Next()
	}
}

func initRoute(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(setDummyID())

	r.GET("/ping", Ping())

	// profile
	r.POST("/profiles", CreateProfile(db))
	r.PUT("/profiles", EditProfile(db))

	return r
}

func TestPong(t *testing.T) {
	// mock db
	db, _, err := getDBMock()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// router
	r := initRoute(db)

	// request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(w, req)

	// expected result
	json := `{"message":"pong"}`

	// assertion
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, json, w.Body.String())
}

func TestCreateProfile(t *testing.T) {
	// mock db
	db, _, err := getDBMock()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// router
	r := initRoute(db)

	// request
	w := httptest.NewRecorder()
	searchID := "meshi"
	name := "meshi reach"
	message := "yoro"
	jsonStr := `{"searchID":"` + searchID + `","name":"` + name + `","message":"` + message + `"}`
	req, _ := http.NewRequest("POST", "/profiles", bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// expected result
	json := `{"status":"success","searchID":"` + searchID + `","name":"` + name + `","message":"` + message + `"}`

	// assertion
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, json, w.Body.String())
}

func TestEditProfile(t *testing.T) {
	// mock db
	db, _, err := getDBMock()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// router
	r := initRoute(db)

	// request
	w := httptest.NewRecorder()
	searchID := "meshiii"
	name := "meshi reach"
	message := "yoro"
	jsonStr := `{"searchID":"` + searchID + `","name":"` + name + `","message":"` + message + `"}`
	req, _ := http.NewRequest("PUT", "/profiles", bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// expected result
	json := `{"status":"success","searchID":"` + searchID + `","name":"` + name + `","message":"` + message + `"}`

	// assertion
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, json, w.Body.String())
}
