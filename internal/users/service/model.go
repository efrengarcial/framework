package service

import (
	base "github.com/efrengarcial/framework/internal/platform/service"
	"time"
)


// User represents a user in the system.
type User struct {
	base.Model
	TenantId		uint64	  `json:"tenantId"`
	Login	  		string    `json:"login" validate:"required" gorm:"not null"`
	Password  		string    `json:"-" validate:"required" gorm:"not null"`
	FirstName 		string    `json:"firstName"`
	LastName  		string    `json:"lastName"`
	Email     		string    `json:"email" validate:"required" gorm:"not null"`
	Activated 		bool	  `json:"activated" validate:"required" gorm:"not null"`
	LangKey   		string    `json:"langKey"`
	ImageUrl  		string    `json:"imageUrl"`
	ActivationKey  	string    `json:"-"`
	ResetKey  		string    `json:"-"`
	ResetDate  		time.Time  ` json:"-"`
	Authorities     []Authority `gorm:"many2many:fw_user_authority;association_autoupdate:false;association_autocreate:false"`
}

type Authority struct {
	base.Model
	Name 		string 			`json:"name" validate:"required" gorm:"not null"`
	TenantId	uint64	  		`json:"tenantId"`
	Privileges  []Privilege `gorm:"many2many:fw_authority_privilege;association_autoupdate:false;association_autocreate:false"`
}


type Privilege struct {
	Name 		string 	`json:"name" gorm:"primary_key"`
}

// Token Entity
type Token struct {
	Token     string    `json:"token" validate:"required"`
	Valid     bool		`json:"valid"`
}

func (user *User) Validate() error {
	return nil
}

func (authority *Authority) Validate() error {
	return nil
}

func (user *User) GetTenantID() uint64 {
	return user.TenantId
}

func (authority *Authority) GetTenantID() uint64 {
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
