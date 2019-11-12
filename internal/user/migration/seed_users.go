package migration

import (
	. "github.com/gobuffalo/nulls"
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

var (
	_roleAdmin              = &Authority{Name: "ROLE_ADMIN"}
	_roleUser               = &Authority{Name: "ROLE_USER"}
	_admin                  = &User{ Login: "_admin" , Email:"_admin@localhost" ,Password: "$2a$10$gSAhZrxMllrbgj/kkK9UceBPpChGWJA7SYIb1Mqo.n5aNLq1/oRrC"}
	_userAuthorityAdmin 	*UserAuthority
	_user               	= &User{ Login: "_user" , Email:"_user@localhost" ,Password: "$2a$10$VEjxo0jq2YG9Rbk2HmX9S.k1uZBGYUHdUcid3g/vfiEl7lwWgOH/K",
								TenantId: NewInt64(1)}
	_userAuthority 			*UserAuthority
)

// SeedUsers inserts the first users
var SeedUsers  = &gormigrate.Migration{
	ID: "SEED_USERS",
	Migrate: func(tx *gorm.DB) error {

		if err := tx.Create(_roleAdmin).Error; err != nil {
			return err
		}
		if err := tx.Create(_roleUser).Error; err != nil {
			return err
		}

		if err :=tx.Create(_admin).Error; err != nil {
			return err
		}
		_userAuthorityAdmin = &UserAuthority{UserID: _admin.ID, AuthorityID: _roleAdmin.ID }
		if err :=tx.Create(_userAuthorityAdmin).Error; err != nil {
			return err
		}

		if err :=tx.Create(&_user).Error; err != nil {
			return err
		}
		_userAuthority = &UserAuthority{UserID: _user.ID, AuthorityID: _roleUser.ID }
		if err :=tx.Create(_userAuthority).Error; err != nil {
			return err
		}

		return nil
	},
	Rollback: func(tx *gorm.DB) error {
		tx.Delete(_userAuthority)
		tx.Delete(_user)
		tx.Delete(_userAuthorityAdmin)
		tx.Delete(_admin)
		tx.Delete(_roleUser)
		tx.Delete(_roleAdmin)

		return nil
	},
}
