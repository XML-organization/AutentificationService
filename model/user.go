package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserCredentials struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email" gorm:"not null;type:string"`
	Password []byte    `json:"password" gorm:"not null;type:string;default:null"`
	Role     Role      `json:"role" gorm:"not null;type:int"`
}

type User struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Role     Role      `json:"role"`
	Name     string    `json:"name"`
	Surname  string    `json:"surname"`
	Address  Address   `json:"address"`
}

func (user *UserCredentials) BeforeCreate(scope *gorm.DB) error {
	user.ID = uuid.New()
	return nil
}
