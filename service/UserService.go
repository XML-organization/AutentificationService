package service

import (
	"autentification_service/model"
	"autentification_service/repository"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo     *repository.UserRepository
	orchestrator *CreateUserOrchestrator
}

func NewUserService(repo *repository.UserRepository, orchestrator *CreateUserOrchestrator) *UserService {
	return &UserService{
		UserRepo:     repo,
		orchestrator: orchestrator,
	}
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

func (service *UserService) Create(user *model.User) error {

	var userCredentials model.UserCredentials
	//hesovanje passworda
	password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)

	userCredentials.ID, _ = uuid.NewUUID()

	println("Ovo je id korisnika koji se treba sacuvati (autentification strana): " + userCredentials.ID.String())

	userCredentials.Password = password
	userCredentials.Email = user.Email
	userCredentials.Role = model.Role(user.Role)

	err := service.UserRepo.CreateUser(&userCredentials)
	if err != nil {
		return err
	}

	err1 := service.orchestrator.Start(user)

	if err1 != nil {
		service.UserRepo.Delete(*user)
		return err1
	}
	return nil
}

func (service *UserService) DeleteUser(user *model.User) error {

	err := service.UserRepo.Delete(*user)
	if err != nil {
		return err
	}
	return nil
}
