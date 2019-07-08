package transport

import (
	"bytes"
	"github.com/efrengarcial/framework/users/pkg/db"
	"github.com/efrengarcial/framework/users/pkg/model"
	"github.com/efrengarcial/framework/users/pkg/repository"
	"github.com/efrengarcial/framework/users/pkg/service"
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
	return  db.Initialize("sqlite3", ":memory:")
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

	repo := repository.NewUserGormRepository(db)
	us := service.NewService(repo, log.With(logger, "component", "users"))
	us = service.NewLoggingService(logger, us)
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

	repo := repository.NewUserGormRepository(db)
	us := service.NewService(repo, log.With(logger, "component", "users"))
	us = service.NewLoggingService(logger, us)
	router := SetupUserRouter(us, logger)

	user := &model.User{Login:"efren.gl" , Email:"efren.garcia@gmail.com" }
	saveUser, err := repo.Insert(user)
	if err != nil {
		t.Fatal(err)
	}

	var  users []model.User
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