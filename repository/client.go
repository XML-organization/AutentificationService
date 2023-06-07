package repository

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetClient(host, user, password, dbname, port string) (*gorm.DB, error) {
	println("host=%s user=%s password=%s dbname=%s port=%s", host, user, password, dbname, port)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
