package repository

import (
	"database/sql"
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/efrengarcial/framework/users/pkg/model"
	"github.com/go-test/deep"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"regexp"
	"strings"
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
	var (
		db  *sql.DB
		err error
	)

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
		`SELECT * FROM "fw_users" WHERE (id = $1)`)).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login"}).
			AddRow(id, login))

	err = repository.Find(user, id)

	require.NoError(t, err)
	require.Nil(t, deep.Equal(&model.User{ Model : model.Model{ ID: id}, Login: login}, user))
}

func Test_repository_Create(t *testing.T) {

	db, mock, err := sqlmock.New()
	defer db.Close()
	require.NoError(t, err)
	DB , err := gorm.Open("postgres", db)
	defer DB.Close()
	require.NoError(t, err)
	DB.LogMode(true)
	repository := NewGormRepository(DB)
	user := &model.User{ Model : model.Model{CreatedAt:  time.Now(),  CreatedBy: "user" , UpdatedAt:  time.Now(), LastModifiedBy: "user"},
			Login:"user", LastName:"user", FirstName:"user", Activated:true , ResetDate:time.Now(), ResetKey: "", LangKey:"us", ActivationKey:"", Email:"user@home",
	ImageUrl:"", Password:"erfsdkkdk"}

	sql := regexp.QuoteMeta(
		`INSERT INTO "fw_users" ("created_at","updated_at","created_by","last_modified_by","login","password","first_name","last_name","email","activated","lang_key","image_url","activation_key","reset_key","reset_date") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15) RETURNING "fw_users"."id" `)
	sql = strings.Replace(sql,"\\", "" ,-1)
	mock.ExpectQuery(sql).
		WithArgs(user.CreatedAt, user.UpdatedAt, user.CreatedBy, user.LastModifiedBy, user.Login, user.Password, user.FirstName,
			user.LastName, user.Email, user.Activated, user.LangKey, user.ImageUrl, user.ActivationKey, user.ResetKey, user.ResetDate ).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(user.ID))


	newUser , err := repository.Insert(user)

	require.NoError(t, err)
	require.NotNil(t, newUser)
}