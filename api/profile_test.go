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

func TestCreateProfileSuccess(t *testing.T) {
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

func TestCreateProfileFailure(t *testing.T) {
	t.Skip("skip failure test")

	// mock db
	db, _, err := getDBMock()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// router
	r := initRoute(db)

	// original data
	searchID := "meshi"
	name := "meshi reach"
	message := "yoro"
	// original := model.User{}
	// original.FirebaseID = "dummy"
	// original.SearchID = searchID
	// original.Name = name
	// original.Message = message

	/* ここでmock dbに書き込まれない */
	// db.Create(&original)

	/* sqlmockのexpect使ってもできない */
	// rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "firebase_id", "search_id", "name", "message"}).AddRow(1, "0000-00-00", "0000-00-00", "0000-00-00", "dummy", searchID, name, message)
	//
	// mock.ExpectQuery(regexp.QuoteMeta(
	// 	`SELECT * FROM "users" ORDER BY "id" LIMIT 1`,
	// )).WillReturnRows(rows)

	// mock.ExpectQuery("^SELECT * FROM users ORDER BY \"id\" LIMIT 1$").WillReturnRows(rows)

	/* dbに書き込まれているかのチェック */
	// search := model.User{}
	// db.First(&search)
	// // db.Where("name = ?", name).Find(&search).Count(&count)
	// fmt.Printf("id=%v\n", search.ID)

	// request
	w := httptest.NewRecorder()
	jsonStr := `{"searchID":"` + searchID + `","name":"` + name + `","message":"` + message + `"}`
	req, _ := http.NewRequest("POST", "/profiles", bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// expected result
	json := `{"status":"failure","error":"This user has already created his/her profile."}`

	// assertion
	assert.Equal(t, http.StatusBadRequest, w.Code)
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
