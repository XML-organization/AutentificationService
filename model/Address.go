package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Address struct {
	ID      uuid.UUID `json:"id"`
	Country string    `json:"country" gorm:"not null;type:string"`
	City    string    `json:"city" gorm:"not null;type:string"`
	Street  string    `json:"street" gorm:"not null;type:string"`
	Number  string    `json:"number" gorm:"not null;type:string"`
}

func (address *Address) BeforeCreate(scope *gorm.DB) error {
	address.ID = uuid.New()
	return nil
}
