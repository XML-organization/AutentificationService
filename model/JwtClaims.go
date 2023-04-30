package model

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type JwtClaims struct {
	Id   uuid.UUID `json:"id" bson:"_id"`
	Role int       `json:"role"`
	jwt.StandardClaims
}
