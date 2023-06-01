package service

import (
	"autentification_service/model"
	"autentification_service/repository"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const SecretKeyForJWT = "v123v1iy2v321sdasada8"

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

	existingUser, _ := service.UserRepo.FindByEmail(user.Email)
	if existingUser.Email == user.Email {
		println("usao u user already exist //////////////////")
		return fmt.Errorf("user already exist")
	}

	var userCredentials model.UserCredentials

	password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)

	userCredentials.ID, _ = uuid.NewUUID()
	userCredentials.Password = password
	userCredentials.Email = user.Email
	userCredentials.Role = model.Role(user.Role)

	err := service.UserRepo.CreateUser(&userCredentials)
	if err != nil {
		return err
	}

	user.ID = userCredentials.ID
	err1 := service.orchestrator.Start(user)

	if err1 != nil {
		service.UserRepo.Delete(*user)
		return err1
	}
	return nil
}

func (service *UserService) ChangePassword(changePassword *model.ChangePasswordDTO) (model.RequestMessage, error) {
	user, err := service.UserRepo.FindByEmail(changePassword.Email)

	if err != nil {
		message := model.RequestMessage{
			Message: "An error occurred, please try again!",
		}
		return message, err
	} else if err := bcrypt.CompareHashAndPassword(user.Password, []byte(changePassword.OldPassword)); err != nil {
		message := model.RequestMessage{
			Message: "The old password is not correct!",
		}
		return message, err
	}

	newPassword, _ := bcrypt.GenerateFromPassword([]byte(changePassword.NewPassword), 14)
	changePassword.NewPassword = string(newPassword)

	return service.UserRepo.ChangePassword(*changePassword)
}

func (service *UserService) ChangeEmail(emails *model.UpdateEmailDTO) error {
	_, err := service.UserRepo.FindByEmail(emails.OldEmail)

	if err != nil {
		return err
	}

	return service.UserRepo.ChangeEmail(*emails)
}

func (service *UserService) DeleteUser(user *model.User) error {

	err := service.UserRepo.Delete(*user)
	if err != nil {
		return err
	}
	return nil
}

func (service *UserService) GetClaimsFrowJwt(tokenString string) (*jwt.Token, model.JwtClaims) {

	claims := &model.JwtClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SecretKeyForJWT), nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, model.JwtClaims{}
		}
		return nil, model.JwtClaims{}
	}

	return token, *claims
}

func (service *UserService) FindByIdUser(id string) (*model.User, error) {
	user, err := service.UserRepo.FindByIdUser(id)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("menu item with id %s not found", id))
	}
	return &user, nil
}

func (service *UserService) DeleteUserCredentials(user *model.UserCredentials) error {

	err := service.UserRepo.DeleteUserCredentials(*user)
	if err != nil {
		return err
	}
	return nil
}
