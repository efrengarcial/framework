package model

import (
	"time"
)

type IModel interface {
	GetId() uint64
	Validate() error
}

type BaseEntity struct {
	Id uint64     `json:"id"`
}

func (base *BaseEntity) GetId() uint64 {
	return base.Id
}

// AbstractAuditingEntity entity
type AbstractAuditingEntity struct {
	CreatedDate     	string
	CreatedBy 			time.Time
	LastModifiedBy     	string
	LastModifiedDate 	time.Time
}


// User Entity
type User struct {
	BaseEntity
	AbstractAuditingEntity  ` json:"-"`
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

// Token Entity
type Token struct {
	Token     string    `json:"token" validate:"required"`
	Valid     bool		`json:"valid"`
}

