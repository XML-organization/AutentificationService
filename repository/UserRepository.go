package repository

import (
	"autentification_service/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	DatabaseConnection *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	err := db.AutoMigrate(&model.UserCredentials{})
	if err != nil {
		return nil
	}

	return &UserRepository{
		DatabaseConnection: db,
	}
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

func (repo *UserRepository) ChangePassword(changePassword model.ChangePasswordDTO) (model.RequestMessage, error) {
	sqlStatementUser := `
		UPDATE user_credentials
		SET password = $2
		WHERE email = $1;`

	dbResult1 := repo.DatabaseConnection.Exec(sqlStatementUser, changePassword.Email, changePassword.NewPassword)

	if dbResult1.Error != nil {
		message := model.RequestMessage{
			Message: "An error occurred, please try again!",
		}
		return message, dbResult1.Error
	}

	message := model.RequestMessage{
		Message: "Success!",
	}
	return message, nil
}

func (repo *UserRepository) ChangeEmail(emails model.UpdateEmailDTO) error {
	sqlStatementUser := `
		UPDATE user_credentials
		SET email = $2
		WHERE email = $1;`

	dbResult1 := repo.DatabaseConnection.Exec(sqlStatementUser, emails.OldEmail, emails.NewEmail)

	if dbResult1.Error != nil {
		return dbResult1.Error
	}
	return nil
}

func (repo *UserRepository) Delete(user model.User) error {
	dbResult := repo.DatabaseConnection.Delete(user)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	println("Rows deleted: ", dbResult.RowsAffected)
	return nil
}

func (repo *UserRepository) FindByIdUser(id string) (model.User, error) {
	user := model.User{}

	dbResult := repo.DatabaseConnection.First(&user, "id = ?", id)

	if dbResult != nil {
		return user, dbResult.Error
	}

	return user, nil
}

func (repo *UserRepository) DeleteUserCredentials(user model.UserCredentials) error {
	dbResult := repo.DatabaseConnection.Delete(user)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	println("Rows deleted: ", dbResult.RowsAffected)
	return nil
}
