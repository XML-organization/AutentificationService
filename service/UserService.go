package service

import (
	"autentification_service/model"
	"autentification_service/repository"
	"fmt"
)

type UserService struct {
	UserRepo *repository.UserRepository
}

func (service *UserService) FindUser(id string) (*model.UserCredentials, error) {
	user, err := service.UserRepo.FindById(id)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("menu item with id %s not found", id))
	}
	return &user, nil
}

func (service *UserService) FindByEmail(email string) (*model.UserCredentials, error) {
	user, err := service.UserRepo.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("User not found!"))
	}

	return &user, nil
}

func (service *UserService) Create(user *model.UserCredentials) error {
	err := service.UserRepo.CreateUser(user)
	if err != nil {
		return err
	}
	return nil
}
