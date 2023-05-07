package service

import (
	"autentification_service/model"

	events "github.com/XML-organization/common/saga/create_user"
	saga "github.com/XML-organization/common/saga/messaging"
)

type CreateUserOrchestrator struct {
	commandPublisher saga.Publisher
	replySubscriber  saga.Subscriber
}

func NewCreateUserOrchestrator(publisher saga.Publisher, subscriber saga.Subscriber) (*CreateUserOrchestrator, error) {
	o := &CreateUserOrchestrator{
		commandPublisher: publisher,
		replySubscriber:  subscriber,
	}
	err := o.replySubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (o *CreateUserOrchestrator) Start(user *model.User) error {
	event := &events.CreateUserCommand{
		Type: events.SaveUser,
		User: *mapUserToSagaUser(user),
	}

	println("Usao sam u metodu start orkestratora")
	return o.commandPublisher.Publish(event)
}

func (o *CreateUserOrchestrator) handle(reply *events.CreateUserReply) {
	command := events.CreateUserCommand{User: reply.User}
	command.Type = o.nextCommandType(reply.Type)
	if command.Type != events.UnknownCommand {
		_ = o.commandPublisher.Publish(command)
	}
}

func (o *CreateUserOrchestrator) nextCommandType(reply events.CreateUserReplyType) events.CreateUserCommandType {
	switch reply {
	case events.UserSaved:
		return events.PrintSuccessful
	case events.UserNotSaved:
		return events.DeleteUserCredentials
	default:
		return events.UnknownCommand
	}
}
