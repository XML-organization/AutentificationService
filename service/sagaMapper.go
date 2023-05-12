package service

import (
	"autentification_service/model"

	events "github.com/XML-organization/common/saga/create_user"
)

func mapUserToSagaUser(u *model.User) *events.User {

	id := " |" + u.ID.String() + " |"
	println("OVO JE ID STRING PRIJE NEGO SE POSLAO USER SERVISU")
	println(id)

	return &events.User{
		Id:       id,
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
