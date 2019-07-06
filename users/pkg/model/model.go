package model

import (
	"time"
)

type IModel interface {
	GetID() uint64
	Validate() error
}

type MultiTenantEntity interface {
	GetTenantID() uint64
}

//BaseModel
type Model struct {
	ID        uint64      	`json:"id,string" gorm:"type:bigint(20) unsigned auto_increment;not null;primary_key`
	CreatedAt *time.Time 	`json:"createdAt,omitempty" gorm:"type:timestamp; not null"`
	UpdatedAt *time.Time 	`json:"updatedAt,omitempty" gorm:"type:timestamp; not null"`
	CreatedBy string 		`json:"createdBy"`
	LastModifiedBy string   `json:"lastModifiedBy"`
	//DeletedAt *time.Time	`json:"deletedAt"`
}

func (base *Model) GetID() uint64 {
	return base.ID
}

// User Entity
type User struct {
	Model
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
	Model
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

type Pageable struct {
	Page    int
	Limit   int
	OrderBy []string
	ShowSQL bool
	Model 	IModel
}
