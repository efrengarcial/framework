package model

import "time"

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
	CreatedAt time.Time 	`json:"createdAt,omitempty" gorm:"type:timestamp; not null"`
	UpdatedAt time.Time 	`json:"updatedAt,omitempty" gorm:"type:timestamp; not null"`
	CreatedBy string 		`json:"createdBy"`
	LastModifiedBy string   `json:"lastModifiedBy"`
	//DeletedAt *time.Time	`json:"deletedAt"`
}

func (base *Model) GetID() uint64 {
	return base.ID
}

type Pageable struct {
	Page    int  `json:"page"`
	Limit   int   `json:"limit"`
	OrderBy []string `json:"orderBy"`
	ShowSQL bool
}

type Pagination struct {
	Count    int         `json:"count"`
	Pages    int         `json:"pages"`
	Records  interface{} `json:"records"`
	Offset   int         `json:"offset"`
	Limit    int         `json:"limit"`
	Page     int         `json:"page"`
	PrevPage int         `json:"prevPage"`
	NextPage int         `json:"nextPage"`
}
