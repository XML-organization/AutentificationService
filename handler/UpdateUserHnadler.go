package handler

import (
	"autentification_service/model"
	"autentification_service/service"
	"log"

	saga "github.com/XML-organization/common/saga/messaging"
	events "github.com/XML-organization/common/saga/update_user"
)

type UpdateUserCommandHandler struct {
	userService       *service.UserService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewUpdateUserCommandHandler(userService *service.UserService, publisher saga.Publisher, subscriber saga.Subscriber) (*UpdateUserCommandHandler, error) {
	o := &UpdateUserCommandHandler{
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

func (handler *UpdateUserCommandHandler) handle(command *events.UpdateUserCommand) {
	reply := events.UpdateUserReply{UpdateUserDTO: command.UpdateUserDTO}

	switch command.Type {
	case events.UpdateUser:
		err := handler.userService.ChangeEmail(&model.UpdateEmailDTO{OldEmail: command.UpdateUserDTO.OldEmail, NewEmail: command.UpdateUserDTO.NewEmail})
		if err != nil {
			reply.Type = events.UserNotUpdated
			log.Println(err)
			log.Println("Saga: User not updated successfuly!")
			break
		}
		log.Println("Saga: Password changed successfuly!")
		reply.Type = events.UserUpdated
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
