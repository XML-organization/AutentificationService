package handler

import (
	"autentification_service/model"

	pb "github.com/XML-organization/common/proto/autentification_service"
)

func MapUserDTOFromLoginRequest(userCredentials *pb.LoginRequest) *model.UserDTO {
	return &model.UserDTO{
		Email:    userCredentials.Email,
		Password: userCredentials.Password,
	}
}
func MapUserFromRegistrationRequest(user *pb.RegistrationRequest) *model.User {
	return &model.User{
		Name:     user.Name,
		Surname:  user.Surname,
		Email:    user.Email,
		Password: user.Password,
		Role:     model.Role(user.Role),
		Country:  user.Country,
		City:     user.City,
		Street:   user.Street,
		Number:   user.Number,
	}
}
