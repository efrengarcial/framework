package transport

import (
	"bytes"
	repository2 "github.com/efrengarcial/framework/pkg/users/repository"
	service2 "github.com/efrengarcial/framework/pkg/users/service"
	"github.com/efrengarcial/framework/pkg/userssers/pkg/util/dbutil"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

func setup() *gorm.DB {
	// Initialize an in-memory database for full integration testing.
	return  dbutil.Initialize("sqlite3", ":memory:")
}

func teardown(db *gorm.DB ) {
	// Closing the connection discards the in-memory database.
	db.Close()
}

func TestCreateHandler(t *testing.T) {
	db := setup()

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = level.NewFilter(logger, level.AllowDebug())
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	repo := repository2.NewUserGormRepository(db)
	us := service2.NewService(repo, log.With(logger, "component", "users"))
	us = service2.NewLoggingService(logger, us)
	router := SetupUserRouter(us, logger)


	//user := &model.User{Login:"efren.gl" , Email:"efren.gl@gmail.com" }
	//repo.Save(user)

	var jsonStr = []byte(`{"login":"efren.gl",  "email" :"efren.gl@gmail.com"}`)
	req, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	// Our API expects a json body, so set the content-type header to make sure it's treated as one.
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	teardown(db)
}

func TestUpdateHandler(t *testing.T) {
	db := setup()

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = level.NewFilter(logger, level.AllowDebug())
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	repo := repository2.NewUserGormRepository(db)
	us := service2.NewService(repo, log.With(logger, "component", "users"))
	us = service2.NewLoggingService(logger, us)
	router := SetupUserRouter(us, logger)

	user := &service2.User{Login: "efren.gl" , Email:"efren.garcia@gmail.com" }
	saveUser, err := repo.Insert(user)
	if err != nil {
		t.Fatal(err)
	}

	var  users []service2.User
	repo.FindAll(&users, "")

	var jsonStr = []byte(`{"id" : "` +  strconv.FormatUint(saveUser.GetID(), 10) + `","login":"efren.gl",  "email" :"efren.gl@gmail.com"}`)
	req, err := http.NewRequest("PUT", "/api/v1/users", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	// Our API expects a json body, so set the content-type header to make sure it's treated as one.
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)


	teardown(db)
}
