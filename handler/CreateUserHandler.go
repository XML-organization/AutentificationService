package handler

import (
	"autentification_service/service"
	"log"

	events "github.com/XML-organization/common/saga/create_user"
	saga "github.com/XML-organization/common/saga/messaging"
)

type CreateUserCommandHandler struct {
	userService       *service.UserService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewCreateUserCommandHandler(userService *service.UserService, publisher saga.Publisher, subscriber saga.Subscriber) (*CreateUserCommandHandler, error) {
	o := &CreateUserCommandHandler{
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

func (handler *CreateUserCommandHandler) handle(command *events.CreateUserCommand) {
	user := mapSagaUserToUser(&command.User)

	reply := events.CreateUserReply{User: command.User}

	switch command.Type {
	case events.PrintSuccessful:
		reply.Type = events.SuccessfulyFinished
	case events.DeleteUserCredentials:
		err := handler.userService.DeleteUser(user)
		if err != nil {
			log.Println(err)
			return
		}
		reply.Type = events.UserCredentialsDeleted
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
