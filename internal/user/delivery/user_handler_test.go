package delivery

import (
	"bytes"
	"github.com/efrengarcial/framework/internal/domain"
	"github.com/efrengarcial/framework/internal/tests"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
	"gopkg.in/gormigrate.v1"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// ProductTests holds methods for each product subtest. This type allows
// passing dependencies for tests while still providing a convenient syntax
// when subtests are registered.
type UserTests struct {
	app       http.Handler
	userToken string
}

func start(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		// you migrations here
	})

	m.InitSchema(func(tx *gorm.DB) error {
		err := tx.AutoMigrate(
			&domain.User{},
			&domain.Authority{},
			&domain.Privilege{},
			// all other tables of you app
		).Error
		if err != nil {
			return err
		}
		user := &domain.User{ Login: "admin" , Email:"admin@example.com" ,Password: "$2a$10$AWLfzWwoq7es.PM3Z1uMieRAeRuck2F.kW9WeEpIdGsk4ykizXLqm"}
		tx.Create(&user)

		// all other foreign keys...
		return nil
	})


	return m.Migrate()
}

//https://medium.com/@yaravind/go-sqlite-on-windows-f91ef2dacfe
func TestCreateHandler(t *testing.T) {

	test := tests.NewIntegration(t)
	defer test.Teardown()
	start(test.DB)

	shutdown := make(chan os.Signal, 1)

	tests := UserTests{
		app: New(shutdown, test.DB,  test.Log, test.Authenticator, test.Cache),
		userToken: test.Token("admin", "Nominable2018."),
	}

	var jsonStr = []byte(`{"login":"efren.gl",  "email" :"efren.gl@gmail.com"}`)
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

/*
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

	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)


	teardown(db)
}
*/
