package handler

import (
	"autentification_service/model"

	changePasswordEvents "github.com/XML-organization/common/saga/change_password"
	createUserEvents "github.com/XML-organization/common/saga/create_user"
)

func mapSagaUserToUser(u *createUserEvents.User) *model.User {

	return &model.User{
		Email:    u.Email,
		Password: u.Password,
		Name:     u.Name,
		Surname:  u.Street,
		Role:     model.Role(u.Role),
		Country:  u.Country,
		City:     u.City,
		Street:   u.Street,
		Number:   u.Number,
	}
}

func mapSagaChangePasswordToChangePasswordDTO(p *changePasswordEvents.ChangePasswordDTO) *model.ChangePasswordDTO {
	return &model.ChangePasswordDTO{
		Email:       p.Email,
		OldPassword: p.OldPassword,
		NewPassword: p.NewPassword,
	}
}
