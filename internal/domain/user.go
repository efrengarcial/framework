package domain

import (
	"time"

	. "github.com/gobuffalo/nulls"
	"github.com/thoas/go-funk"
)

type Password string

// User represents a user in the system.
type User struct {
	Model
	TenantId      Int64  `json:"tenantId"`
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
	Authorities   []Authority `gorm:"many2many:fw_user_authority;association_autoupdate:false;association_autocreate:false"`
	Permissions   []string    `json:"permissions" gorm:"-"` // Ignore this field
}

func (Password) MarshalJSON() ([]byte, error) {
	return []byte(`""`), nil
}

type Authority struct {
	Model
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
	UserName   string `json:"username"  binding:"required"`
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

func (user *User) ConfigPermissions() {
	if len(user.Authorities) > 0  {
		r := funk.Map(user.Authorities, func(a Authority) string {
			return a.Name
		})
		user.Permissions =  append( user.Permissions , r.([]string)...)
		for _, authority := range user.Authorities {
			if len(authority.Privileges) > 0 {
				p := funk.Map(authority.Privileges, func(p Privilege) string {
					return p.Name
				})
				user.Permissions = append(user.Permissions, p.([]string)...)
			}
		}
		user.Authorities = nil

	}
}

func (user *User) GetPermissions() []string {
	return user.Permissions
}

// HasRole returns true if the user has at least one of the provided permission.
func (user *User) HasPermission(permissions ...string) bool {

	for _, has := range user.GetPermissions() {
		for _, want := range permissions {
			if has == want {
				return true
			}
		}
	}

	return false
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
