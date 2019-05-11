package repository

import (
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/efrengarcial/framework/users/pkg/model"
	"github.com/go-test/deep"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"log"
	"regexp"
	"testing"
	"time"
)
type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}


func  Test_repository_Find(t *testing.T) {

	db, mock, err := sqlmock.New()
	defer db.Close()
	require.NoError(t, err)
	DB , err := gorm.Open("postgres", db)
	defer DB.Close()
	require.NoError(t, err)
	DB.LogMode(true)
	repository := NewGormRepository(DB)
	user := new(model.User)

	var (
		id  uint64  = 1
		login = "user"
	)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "fw_user" WHERE (id = $1)`)).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login"}).
			AddRow(id, login))

	err = repository.Find(user, id)

	require.NoError(t, err)
	require.Nil(t, deep.Equal(&model.User{ Model : model.Model{ ID: id}, Login: login}, user))
}

func Test_repository_Create(t *testing.T) {

	var id  uint64  = 1
	db, mock, err := sqlmock.New()
	defer db.Close()
	require.NoError(t, err)
	DB , err := gorm.Open("postgres", db)
	defer DB.Close()
	require.NoError(t, err)
	DB.LogMode(true)
	repository := NewGormRepository(DB)

	user := &model.User{ Model : model.Model{LastModifiedBy: "user"}, Login:"user", LastName:"user", FirstName:"user",
		Activated:true , ResetKey: "", LangKey:"us", ActivationKey:"", Email:"user@home", ImageUrl:"", Password:"erfsdkkdk" ,
		Authorities: []model.Authority{ {Name: "ROLE_ADMIN"} }}

	mock.ExpectQuery( regexp.QuoteMeta(
		`INSERT INTO "fw_user" `)).
		WithArgs(AnyTime{}, AnyTime{}, user.CreatedBy, user.LastModifiedBy, user.Login, user.Password, user.FirstName,
			user.LastName, user.Email, user.Activated, user.LangKey, user.ImageUrl, user.ActivationKey, user.ResetKey, AnyTime{}).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(id))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "fw_authority" `)).
		WithArgs("ROLE_ADMIN").
		WillReturnRows(sqlmock.NewRows([]string{"name"}).
			AddRow("ROLE_ADMIN"))

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "fw_user_authority"  `)).
		WithArgs(id, "ROLE_ADMIN", id, "ROLE_ADMIN").WillReturnResult(sqlmock.NewResult(1, 1))



	newUser , err := repository.Insert(user)
	log.Println(newUser.GetID())

	require.NoError(t, err)
	require.NotNil(t, newUser)
}

func Test_repository_Save(t *testing.T) {

	var id  uint64  = 1
	db, mock, err := sqlmock.New()
	defer db.Close()
	require.NoError(t, err)
	DB , err := gorm.Open("postgres", db)
	defer DB.Close()
	require.NoError(t, err)
	DB.LogMode(true)
	repository := NewGormRepository(DB)
	user := &model.User{ Model : model.Model{LastModifiedBy: "user"}, Login:"user", LastName:"user", FirstName:"user",
		Activated:true , ResetKey: "", LangKey:"us", ActivationKey:"", Email:"user@home", ImageUrl:"", Password:"erfsdkkdk"}

	sql := regexp.QuoteMeta(
		`INSERT INTO "fw_user" `)
	mock.ExpectQuery(sql).
		WithArgs(AnyTime{}, AnyTime{}, user.CreatedBy, user.LastModifiedBy, user.Login, user.Password, user.FirstName,
			user.LastName, user.Email, user.Activated, user.LangKey, user.ImageUrl, user.ActivationKey, user.ResetKey, AnyTime{}).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(id))


	actualId , err := repository.Save(user)
	log.Println(id)

	require.NoError(t, err)
	require.Equal(t, uint64(id),actualId)
}