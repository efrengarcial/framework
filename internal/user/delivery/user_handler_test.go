package delivery

import (
	"bytes"
	"context"
	"github.com/efrengarcial/framework/internal/domain"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/efrengarcial/framework/internal/platform/database"
	"github.com/efrengarcial/framework/internal/user/repository"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func setup() *gorm.DB {
	// Initialize an in-memory database for full integration testing.
	db := database.Initialize("sqlite3", ":memory:")
	db.AutoMigrate(&domain.User{}, &domain.Authority{}, &domain.Privilege{})

	return  db
}

func teardown(db *gorm.DB ) {
	// Closing the connection discards the in-memory database.
	db.Close()
}

func TestCreateHandler(t *testing.T) {
	db := setup()

	logger := log.New()
	logger.Out = os.Stdout
	logger.Level = log.InfoLevel
	logger.Formatter = &log.JSONFormatter{}
	shutdown := make(chan os.Signal, 1)
	server := New(shutdown, db, logger)

	//user := &model.User{UserName:"efren.gl" , Email:"efren.gl@gmail.com" }
	//repo.Save(user)

	var jsonStr = []byte(`{"login":"efren.gl",  "email" :"efren.gl@gmail.com"}`)
	req, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	// Our API expects a json body, so set the content-type header to make sure it's treated as one.
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	teardown(db)
}

func TestUpdateHandler(t *testing.T) {
	db := setup()

	logger := log.New()
	logger.Out = os.Stdout
	logger.Level = log.InfoLevel
	logger.Formatter = &log.JSONFormatter{}

	repo := repository.NewUserGormRepository(db)
	shutdown := make(chan os.Signal, 1)
	server := New(shutdown, db, logger)
	user := &domain.User{Login: "efren.gl" , Email:"efren.garcia@gmail.com" }
	saveUser, err := repo.Insert(context.Background(), user)
	if err != nil {
		t.Fatal(err)
	}

	var  users []domain.User
	repo.FindAll(&users, "")

	var jsonStr = []byte(`{"id" : "` +  strconv.FormatUint(saveUser.GetID(), 10) + `","login":"efren.gl",  "email" :"efren.gl@gmail.com"}`)
	req, err := http.NewRequest("PUT", "/api/v1/users", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	// Our API expects a json body, so set the content-type header to make sure it's treated as one.
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)


	teardown(db)
}
