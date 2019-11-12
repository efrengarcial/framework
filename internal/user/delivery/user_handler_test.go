package delivery

import (
	"bytes"
	"context"
	"github.com/efrengarcial/framework/internal/domain"
	"github.com/efrengarcial/framework/internal/tests"
	"github.com/efrengarcial/framework/internal/user/migration"
	"github.com/efrengarcial/framework/internal/user/repository"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

// ProductTests holds methods for each product subtest. This type allows
// passing dependencies for tests while still providing a convenient syntax
// when subtests are registered.
type UserTests struct {
	app       http.Handler
	userToken string
}

//https://medium.com/@yaravind/go-sqlite-on-windows-f91ef2dacfe
func TestCreateHandler(t *testing.T) {

	test := tests.NewIntegration(t)
	defer test.Teardown()
	migration.Start(test.DB)

	shutdown := make(chan os.Signal, 1)

	tests := UserTests{
		app: New(shutdown, test.DB,  test.Log, test.Authenticator, test.Cache),
		userToken: test.Token("admin", "admin"),
	}

	var jsonStr = []byte(`{"login":"efren.gl",  "email" :"efren.gl@gmail.com", "password": "pegasso" }`)
	req, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	// Our API expects a json body, so set the content-type header to make sure it's treated as one.
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+tests.userToken)

	w := httptest.NewRecorder()
	tests.app.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}


func TestUpdateHandler(t *testing.T) {
	test := tests.NewIntegration(t)
	defer test.Teardown()
	migration.Start(test.DB)

	shutdown := make(chan os.Signal, 1)

	tests := UserTests{
		app: New(shutdown, test.DB,  test.Log, test.Authenticator, test.Cache),
		userToken: test.Token("admin", "admin"),
	}

	repo := repository.NewUserGormRepository(test.DB)
	user := &domain.User{Login: "efren.gl" , Email:"efren.garcia@gmail.com" }
	err := repo.Insert(context.Background(), user)
	if err != nil {
		t.Fatal(err)
	}

	var  users []domain.User
	repo.FindAll(&users, "")

	var jsonStr = []byte(`{"id" : "` +  strconv.FormatUint(user.GetID(), 10) + `","login":"efren.gl",  "email" :"efren.gl@gmail.com"}`)
	req, err := http.NewRequest("PUT", "/api/v1/users", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	// Our API expects a json body, so set the content-type header to make sure it's treated as one.
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+tests.userToken)

	w := httptest.NewRecorder()
	tests.app.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

