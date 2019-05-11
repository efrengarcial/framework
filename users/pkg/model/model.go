package model

import (
	"time"
)

type IModel interface {
	GetID() uint64
	Validate() error
}

type Model struct {
	ID        uint64      	`json:"id" gorm:"type:bigserial;primary_key"`
	CreatedAt time.Time 	`json:"createdAt" gorm:"type:timestamp"`
	UpdatedAt time.Time 	`json:"updatedAt" gorm:"type:timestamp"`
	CreatedBy string 		`json:"createdBy"`
	LastModifiedBy string    `json:"lastModifiedBy"`
}

func (base *Model) GetID() uint64 {
	return base.ID
}

// User Entity
type User struct {
	Model
	Login	  		string    `json:"login" validate:"required"`
	Password  		string    `json:"password" validate:"required"`
	FirstName 		string    `json:"firstName"`
	LastName  		string    `json:"lastName"`
	Email     		string    `json:"email" validate:"required"`
	Activated 		bool		`json:"activated" validate:"required"`
	LangKey   		string    	`json:"langKey"`
	ImageUrl  		string    `json:"imageUrl"`
	ActivationKey  	string    ` json:"-"`
	ResetKey  		string    ` json:"-"`
	ResetDate  		time.Time    ` json:"-"`
}

func (user *User) Validate() error {
	return nil
}

// Set User's table name to be `fw_users`
func (User) TableName() string {
	return "fw_users"
}

// Token Entity
type Token struct {
	Token     string    `json:"token" validate:"required"`
	Valid     bool		`json:"valid"`
}

