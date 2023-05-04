package model

import (
	"github.com/dgrijalva/jwt-go"
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

type JwtClaims struct {
	Id   uuid.UUID `json:"id" bson:"_id"`
	Role int       `json:"role"`
	jwt.StandardClaims
}

type RequestMessage struct {
	Message string `json:"message"`
}

type Role int

const (
	Host Role = iota
	Guest
	NK
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
	Country  string    `json:"country" gorm:"not null;type:string"`
	City     string    `json:"city" gorm:"not null;type:string"`
	Street   string    `json:"street" gorm:"not null;type:string"`
	Number   string    `json:"number" gorm:"not null;type:string"`
}

func (user *UserCredentials) BeforeCreate(scope *gorm.DB) error {
	user.ID = uuid.New()
	return nil
}

type UserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
