package service

import (
	"time"

	base "github.com/efrengarcial/framework/internal/platform/model"
	. "github.com/markbates/pop/nulls"
	"github.com/thoas/go-funk"
)


// User represents a user in the system.
type User struct {
	base.Model
	TenantId		Int64     `json:"tenantId"`
	Login	  		string    `json:"login" binding:"required" gorm:"not null"`
	Password  		string    `json:"password" binding:"required" gorm:"not null"`
	FirstName 		string    `json:"firstName"`
	LastName  		string    `json:"lastName"`
	Email     		string    `json:"email" binding:"required" gorm:"not null"`
	Activated 		bool	  `json:"activated" gorm:"not null"`
	LangKey   		string    `json:"langKey"`
	ImageUrl  		string    `json:"imageUrl"`
	ActivationKey  	string    `json:"-"`
	ResetKey  		string    `json:"-"`
	ResetDate  		time.Time  ` json:"-"`
	Authorities     []Authority `gorm:"many2many:fw_user_authority;association_autoupdate:false;association_autocreate:false"`
}

type Authority struct {
	base.Model
	Name 		string 			`json:"name" binding:"required" gorm:"not null"`
	TenantId	Int64	  		`json:"tenantId"`
	Privileges  []Privilege 	`gorm:"many2many:fw_authority_privilege;association_autoupdate:false;association_autocreate:false"`
}


type Privilege struct {
	Name 		string 	`json:"name" gorm:"primary_key"`
}

// Token Entity
type Token struct {
	Token     string    `json:"token"`
	Valid     bool		`json:"valid"`
}


type LoginVM struct {
	Login      string `json:"login"  binding:"required"`
	Password   string `json:"password"  binding:"required"`
	RememberMe bool   `json:"rememberMe"`
}

func (user *User) Validate() error {
	return nil
}

func (authority *Authority) Validate() error {
	return nil
}

func (user *User) GetTenantID() Int64 {
	return user.TenantId
}

func (user *User) GetRoles() []string {
	if len(user.Authorities) > 0  {
		r := funk.Map(user.Authorities, func(a Authority) string {
			return a.Name
		})
		return r.([]string)
	}
	return nil
}

func (authority *Authority) GetTenantID() Int64 {
	return authority.TenantId
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
func (Privilege) TableName() string {
	return "fw_privilege"
}
