package handler

import (
	"autentification_service/service"
	"log"

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
		log.Println(err)
		return nil, err
	}
	return o, nil
}

func (handler *ChangePasswordCommandHandler) handle(command *events.ChangePasswordCommand) {
	reply := events.ChangePasswordReply{ChangePasswordDTO: command.ChagePasswordDTO}

	log.Println("Change password: Usao sam u handle metodu na autentification strani")
	log.Println("Ovo je tip comande koju sam dobio:", command.Type)

	switch command.Type {
	case events.ChangePassword:
		log.Println("Novi password" + command.ChagePasswordDTO.NewPassword)
		log.Println("Stari password" + command.ChagePasswordDTO.OldPassword)
		_, err := handler.userService.ChangePassword(mapSagaChangePasswordToChangePasswordDTO(&command.ChagePasswordDTO))
		if err != nil {
			reply.Type = events.PasswordNotChanged
			log.Println(err)
			log.Println("Saga: User password dont changed successfuly!")
			break
		}
		log.Println("Saga: Password changed successfuly!")
		reply.Type = events.PasswordChanged
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
