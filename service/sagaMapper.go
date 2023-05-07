package service

import (
	"autentification_service/model"

	events "github.com/XML-organization/common/saga/create_user"
)

func mapUserToSagaUser(u *model.User) *events.User {
	return &events.User{
		ID:       u.ID.String(),
		Name:     u.Name,
		Surname:  u.Surname,
		Email:    u.Email,
		Password: u.Password,
		Role:     events.Role(u.Role),
		Country:  u.Country,
		City:     u.City,
		Street:   u.Street,
		Number:   u.Number,
	}
}
