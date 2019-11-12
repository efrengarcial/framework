package migration

import (
	"github.com/efrengarcial/framework/internal/domain"
	. "github.com/gobuffalo/nulls"
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
	"time"
)

type Password string

type User struct {
	domain.Model
	TenantId      Int64  	  `json:"tenantId"`
	Login         string      `json:"login" binding:"required" gorm:"not null"`
	Password      Password    `json:"password"`
	FirstName     string      `json:"firstName"`
	LastName      string      `json:"lastName"`
	Email         string      `json:"email" binding:"required" gorm:"not null"`
	Activated     bool        `json:"activated" gorm:"not null"`
	LangKey       string      `json:"langKey"`
	ImageUrl      string      `json:"imageUrl"`
	ActivationKey string      `json:"-"`
	ResetKey      string      `json:"-"`
	ResetDate     time.Time   `json:"resetDate"`
}

type Authority struct {
	domain.Model
	Name 		string 			`json:"name" binding:"required" gorm:"not null"`
	TenantId	Int64	  		`json:"tenantId"`
}

type UserAuthority struct {
	UserID uint64 `gorm:"primary_key;auto_increment:false"`
	AuthorityID uint64 `gorm:"primary_key;auto_increment:false"`
}

type Privilege struct {
	Name 		string 	`json:"name" gorm:"primary_key"`
}

type AuthorityPrivilege struct {
	AuthorityID uint64 `gorm:"primary_key;auto_increment:false"`
	PrivilegeName string `gorm:"primary_key;auto_increment:false"`
}


// Set User's table name to be `fw_user`
func (User) TableName() string {
	return "fw_user"
}

// Set User's table name to be `fw_authority`
func (Authority) TableName() string {
	return "fw_authority"
}

// Set User's table name to be `fw_authority`
func (UserAuthority) TableName() string {
	return "fw_user_authority"
}

// Set User's table name to be `fw_authority`
func (Privilege) TableName() string {
	return "fw_privilege"
}

func (AuthorityPrivilege) TableName() string {
	return "fw_authority_privilege"
}


func Start(db *gorm.DB) error {

	m := gormigrate.New(db, gormigrate.DefaultOptions, nil)


	m.InitSchema(func(tx *gorm.DB) error {
		err := tx.AutoMigrate(
			&User{},
			&Authority{},
			&UserAuthority{},
			&Privilege{},
			&AuthorityPrivilege{},
			// all other tables of you app
		).Error
		if err != nil {
			return err
		}

		if tx.Dialect().GetName() != "sqlite3" {
			if err :=tx.Model(&UserAuthority{}).AddForeignKey("user_id", " fw_user(id)", "RESTRICT", "RESTRICT").Error; err != nil {
				return err
			}
			if err :=tx.Model(&UserAuthority{}).AddForeignKey("authority_id", " fw_authority(id)", "RESTRICT", "RESTRICT").Error; err != nil {
				return err
			}
			if err :=tx.Model(&AuthorityPrivilege{}).AddForeignKey("authority_id", " fw_authority(id)", "RESTRICT", "RESTRICT").Error; err != nil {
				return err
			}
			if err :=tx.Model(&AuthorityPrivilege{}).AddForeignKey("privilege_name", " fw_privilege(name)", "RESTRICT", "RESTRICT").Error; err != nil {
				return err
			}
		}

		if err :=tx.Model(&User{}).AddUniqueIndex("idx_user_tenant_id_login", "tenant_id", "login").Error; err != nil {
			return err
		}

		// all other foreign keys...
		return nil
	})

	if err := m.Migrate(); err != nil {
		return err
	}

	m = gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		SeedUsers,
	})

	return m.Migrate()
}
