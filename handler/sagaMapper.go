package handler

import (
	"autentification_service/model"

	events "github.com/XML-organization/common/saga/create_user"
	"github.com/google/uuid"
)

func mapSagaUserToUser(u *events.User) *model.User {

	id, _ := uuid.Parse(u.ID)

	return &model.User{
		ID:       id,
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
