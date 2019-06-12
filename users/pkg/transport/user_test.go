package transport

import (
	"bytes"
	"github.com/efrengarcial/framework/users/pkg/db"
	"github.com/efrengarcial/framework/users/pkg/repository"
	"github.com/efrengarcial/framework/users/pkg/service"
	"github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"net/http"
	"net/http/httptest"
	"os"
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
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	repo := repository.NewGormRepository(db)
	us := service.NewService(repo, log.With(logger, "component", "users"))
	handler := userHandler{us, logger}

	var jsonStr = []byte(`{"id":"1"}`)
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	// Our API expects a json body, so set the content-type header to make sure it's treated as one.
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	rr := httptest.NewRecorder()

	http.HandlerFunc(handler.createUser).ServeHTTP(rr, req)

	// Test that the status code is correct.
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Status code is invalid. Expected %d. Got %d instead", http.StatusCreated, status)
	}


	teardown(db)
}
