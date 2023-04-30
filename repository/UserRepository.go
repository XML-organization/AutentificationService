package repository

import (
	"autentification_service/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	DatabaseConnection *gorm.DB
}

func (repo *UserRepository) FindById(id string) (model.UserCredentials, error) {
	user := model.UserCredentials{}

	dbResult := repo.DatabaseConnection.First(&user, "id = ?", id)

	if dbResult != nil {
		return user, dbResult.Error
	}

	return user, nil
}

func (repo *UserRepository) FindByEmail(email string) (model.UserCredentials, error) {
	user := model.UserCredentials{}

	dbResult := repo.DatabaseConnection.First(&user, "email = ?", email)

	if dbResult != nil {
		return user, dbResult.Error
	}

	return user, nil
}

func (repo *UserRepository) CreateUser(user *model.UserCredentials) error {
	dbResult := repo.DatabaseConnection.Create(user)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	println("Rows affected: ", dbResult.RowsAffected)
	return nil
}
