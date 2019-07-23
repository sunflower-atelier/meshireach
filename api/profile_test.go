package api

import (
	"bytes"
	"log"
	"meshireach/db/model"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/assert"
)

// not used now
func initDB() *gorm.DB {
	db, err := gorm.Open("mysql", "root:@tcp(db:3306)/meshireach_test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Friend_Relation{})
	db.AutoMigrate(&model.Event{})
	db.AutoMigrate(&model.Participants{})

	return db
}

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

var db *gorm.DB

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("tcp://192.168.99.100:2376")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("mysql", "8.0.16", nil)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		db, err = gorm.Open("mysql", "root:@tcp(db:3306)/meshireach_test?charset=utf8&parseTime=True&loc=Local")
		if err != nil {
			return err
		}
		return err
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestPong(t *testing.T) {
	// test db
	// db := initDB()
	// defer db.Close()

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
	// test db
	// db := initDB()
	// defer db.Close()
	// defer db.Delete(model.User{})

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

	/* dbに書き込まれているかのチェック */
	search := model.User{}
	db.First(&search)
	t.Log("---name---")
	t.Log(search.Name)
}

func TestCreateProfileFailure(t *testing.T) {
	t.Skip("skip failure test")

	// test db
	// db := initDB()
	// defer db.Close()
	// defer db.Delete(model.User{})

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
	// test db
	// db := initDB()
	// defer db.Close()
	// defer db.Delete(model.User{})

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
