package repository

import (
	"context"
	"database/sql/driver"
	"github.com/efrengarcial/framework/internal/domain"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/efrengarcial/framework/internal/platform/repository"
	"github.com/efrengarcial/framework/internal/users"
	"github.com/go-test/deep"
	"github.com/jinzhu/gorm"
	. "github.com/markbates/pop/nulls"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	repository := repository.NewGormRepository(DB)
	user := new(users.User)

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

	assert.NoError(t, err)
	assert.Nil(t, deep.Equal(&users.User{ Model : domain.Model{ ID: id}, Login: login}, user))
}

func Test_repository_Create(t *testing.T) {

	var ( id  uint64  = 1
		 tenantId = NewInt64(10)
		 )
	db, mock, err := sqlmock.New()
	defer db.Close()
	require.NoError(t, err)
	DB , err := gorm.Open("postgres", db)
	defer DB.Close()
	require.NoError(t, err)
	DB.LogMode(true)
	repository := repository.NewGormRepository(DB)
	user := &users.User{ Model : domain.Model{CreatedBy: "user", LastModifiedBy: "user"}, TenantId: tenantId, Login:"user", LastName:"user", FirstName:"user",
		Activated:true , ResetKey: "", LangKey:"us", ActivationKey:"", Email:"user@home", ImageUrl:"", Password:"erfsdkkdk"}

	mock.ExpectQuery( regexp.QuoteMeta(
		`INSERT INTO "fw_user" `)).
		WithArgs(AnyTime{}, AnyTime{}, user.CreatedBy, user.LastModifiedBy, user.TenantId, user.Login, user.Password, user.FirstName,
			user.LastName, user.Email, user.Activated, user.LangKey, user.ImageUrl, user.ActivationKey, user.ResetKey, AnyTime{}).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(id))

	newUser , err := repository.Insert(context.Background(), user)
	log.Println(newUser.GetID())

	assert.NoError(t, err)
	assert.NotNil(t, newUser)
}

func Test_repository_Create_ExistingAuthority(t *testing.T) {

	var ( id  uint64  = 1
		tenantId = NewInt64(10)
		roleAdmin = "ROLE_ADMIN"
	)
	db, mock, err := sqlmock.New()
	defer db.Close()
	require.NoError(t, err)
	DB , err := gorm.Open("postgres", db)
	defer DB.Close()
	require.NoError(t, err)
	DB.LogMode(true)
	repository := repository.NewGormRepository(DB)

	authority := users.Authority{Model : domain.Model{CreatedBy: "user", LastModifiedBy: "user"}, Name: roleAdmin , TenantId:tenantId}

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "fw_authority"  `)).
		WithArgs(AnyTime{}, AnyTime{}, authority.CreatedBy, authority.LastModifiedBy, roleAdmin, tenantId).WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(id))

	repository.Insert(context.Background(), &authority)

	existAuthority :=  users.Authority{Model : domain.Model{ID: id}}

	user := &users.User{ Model : domain.Model{CreatedBy: "user", LastModifiedBy: "user"}, TenantId: tenantId, Login:"user", LastName:"user", FirstName:"user",
		Activated:true , ResetKey: "", LangKey:"us", ActivationKey:"", Email:"user@home", ImageUrl:"", Password:"erfsdkkdk" ,
		Authorities: []users.Authority{existAuthority }}

	mock.ExpectQuery( regexp.QuoteMeta(
		`INSERT INTO "fw_user" `)).
		WithArgs(AnyTime{}, AnyTime{}, user.CreatedBy, user.LastModifiedBy, user.TenantId, user.Login, user.Password, user.FirstName,
			user.LastName, user.Email, user.Activated, user.LangKey, user.ImageUrl, user.ActivationKey, user.ResetKey, AnyTime{}).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(id))



	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "fw_user_authority"  `)).
		WithArgs(id, id, id, id).WillReturnResult(sqlmock.NewResult(1, 1))

	newUser , err := repository.Insert(context.Background(), user)
	log.Println(newUser.GetID())

	assert.NoError(t, err)
	assert.NotNil(t, newUser)
}

func Test_repository_Save(t *testing.T) {

	var ( id  uint64  = 1
		tenantId = NewInt64(10)
	)
	db, mock, err := sqlmock.New()
	defer db.Close()
	require.NoError(t, err)
	DB , err := gorm.Open("postgres", db)
	defer DB.Close()
	require.NoError(t, err)
	DB.LogMode(true)
	repository := repository.NewGormRepository(DB)
	user := &users.User{ Model : domain.Model{CreatedBy: "user", LastModifiedBy: "user"}, TenantId: tenantId,  Login:"user", LastName:"user", FirstName:"user",
		Activated:true , ResetKey: "", LangKey:"us", ActivationKey:"", Email:"user@home", ImageUrl:"", Password:"erfsdkkdk"}

	sql := regexp.QuoteMeta(
		`INSERT INTO "fw_user" `)
	mock.ExpectQuery(sql).
		WithArgs(AnyTime{}, AnyTime{}, user.CreatedBy, user.LastModifiedBy, user.TenantId , user.Login, user.Password, user.FirstName,
			user.LastName, user.Email, user.Activated, user.LangKey, user.ImageUrl, user.ActivationKey, user.ResetKey, AnyTime{}).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(id))


	actualId , err := repository.Save(user)
	log.Println(id)

	assert.NoError(t, err)
	assert.Equal(t, uint64(id),actualId)
}

func  Test_repository_FindAll(t *testing.T) {

	db, mock, err := sqlmock.New()
	defer db.Close()
	require.NoError(t, err)
	DB , err := gorm.Open("postgres", db)
	defer DB.Close()
	require.NoError(t, err)
	DB.LogMode(true)
	repository := repository.NewGormRepository(DB)

	var (
		id  uint64  = 1
		login = "user"
		users []users.User
	)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "fw_user"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login"}).
			AddRow(id, login))

	err = repository.FindAll(&users, "id > 0 " )

	assert.NoError(t, err)
	assert.Len(t, users, 1)
}

func  Test_repository_FindAllPageable(t *testing.T) {

	db, mock, err := sqlmock.New()
	defer db.Close()
	require.NoError(t, err)
	DB , err := gorm.Open("postgres", db)
	defer DB.Close()
	require.NoError(t, err)
	DB.LogMode(true)
	repository := repository.NewGormRepository(DB)

	var (
		id  uint64  = 1
		login = "user"
		users []users.User
	)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "fw_user"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login"}).
			AddRow(id, login))

	mock.ExpectQuery(regexp.QuoteMeta(
		` SELECT count(*) FROM "fw_user"`)).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).
			AddRow(1))

	pageable := domain.Pageable{Page: 1 , Limit: 10 , OrderBy: []string{"id desc"}}
	_, err  = repository.FindAllPageable(context.Background(), &pageable, &users, "id > 0 ")

	assert.NoError(t, err)
	assert.Len(t, users, 1)
}
