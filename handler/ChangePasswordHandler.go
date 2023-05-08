package handler

import (
	"autentification_service/service"

	events "github.com/XML-organization/common/saga/change_password"
	saga "github.com/XML-organization/common/saga/messaging"
)

type ChangePasswordCommandHandler struct {
	userService       *service.UserService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewChangePasswordCommandHandler(userService *service.UserService, publisher saga.Publisher, subscriber saga.Subscriber) (*ChangePasswordCommandHandler, error) {
	o := &ChangePasswordCommandHandler{
		userService:       userService,
		replyPublisher:    publisher,
		commandSubscriber: subscriber,
	}
	err := o.commandSubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (handler *ChangePasswordCommandHandler) handle(command *events.ChangePasswordCommand) {
	reply := events.ChangePasswordReply{ChangePasswordDTO: command.ChagePasswordDTO}

	println("Change password: Usao sam u handle metodu na autentification strani")
	println("Ovo je tip comande koju sam dobio: %v", command.Type)

	switch command.Type {
	case events.ChangePassword:
		println("Novi password" + command.ChagePasswordDTO.NewPassword)
		println("Stari password" + command.ChagePasswordDTO.OldPassword)
		_, err := handler.userService.ChangePassword(mapSagaChangePasswordToChangePasswordDTO(&command.ChagePasswordDTO))
		if err != nil {
			reply.Type = events.PasswordNotChanged
			println("Saga: User password dont changed successfuly!")
			break
		}
		println("Saga: Password changed successfuly!")
		reply.Type = events.PasswordChanged
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
