package startup

import (
	"autentification_service/handler"
	"autentification_service/repository"
	"autentification_service/service"
	"autentification_service/startup/config"
	"fmt"
	"log"
	"net"

	pb "github.com/XML-organization/common/proto/autentification_service"
	saga "github.com/XML-organization/common/saga/messaging"
	"github.com/XML-organization/common/saga/messaging/nats"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type Server struct {
	config *config.Config
}

func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}

const (
	QueueGroup = "autentification_service"
)

func (server *Server) Start() {
	postgresClient := server.initPostgresClient()
	userRepo := server.initUserRepository(postgresClient)

	//create user
	commandPublisher := server.initPublisher(server.config.CreateUserCommandSubject)
	replySubscriber := server.initSubscriber(server.config.CreateUserReplySubject, QueueGroup)
	createUserOrchestrator := server.initCreateOrderOrchestrator(commandPublisher, replySubscriber)

	userService := server.initUserService(userRepo, createUserOrchestrator)
	userHandler := server.initUserHandler(userService)

	//create user
	commandSubscriber1 := server.initSubscriber(server.config.CreateUserCommandSubject, QueueGroup)
	replyPublisher1 := server.initPublisher(server.config.CreateUserReplySubject)
	server.initCreateUserHandler(userService, replyPublisher1, commandSubscriber1)

	//change password
	commandSubscriber2 := server.initSubscriber(server.config.ChangePasswordCommandSubject, QueueGroup)
	replyPublisher2 := server.initPublisher(server.config.ChangePasswordReplySubject)
	server.initChangePasswordHandler(userService, replyPublisher2, commandSubscriber2)

	server.startGrpcServer(userHandler)
}

func (server *Server) initCreateUserHandler(service *service.UserService, publisher saga.Publisher, subscriber saga.Subscriber) {
	_, err := handler.NewCreateUserCommandHandler(service, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}

func (server *Server) initChangePasswordHandler(service *service.UserService, publisher saga.Publisher, subscriber saga.Subscriber) {
	_, err := handler.NewChangePasswordCommandHandler(service, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}

func (server *Server) initPublisher(subject string) saga.Publisher {
	publisher, err := nats.NewNATSPublisher(
		server.config.NatsHost, server.config.NatsPort,
		server.config.NatsUser, server.config.NatsPass, subject)
	if err != nil {
		log.Fatal(err)
	}
	return publisher
}

func (server *Server) initSubscriber(subject, queueGroup string) saga.Subscriber {
	subscriber, err := nats.NewNATSSubscriber(
		server.config.NatsHost, server.config.NatsPort,
		server.config.NatsUser, server.config.NatsPass, subject, queueGroup)
	if err != nil {
		log.Fatal(err)
	}
	return subscriber
}

func (server *Server) initCreateOrderOrchestrator(publisher saga.Publisher, subscriber saga.Subscriber) *service.CreateUserOrchestrator {
	orchestrator, err := service.NewCreateUserOrchestrator(publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
	return orchestrator
}

func (server *Server) initPostgresClient() *gorm.DB {
	client, err := repository.GetClient(
		server.config.AutentificationDBHost, server.config.AutentificationDBUser,
		server.config.AutentificationDBPass, server.config.AutentificationDBName,
		server.config.AutentificationDBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) initUserRepository(client *gorm.DB) *repository.UserRepository {
	return repository.NewUserRepository(client)
}

func (server *Server) initUserService(repo *repository.UserRepository, orchestrator *service.CreateUserOrchestrator) *service.UserService {
	return service.NewUserService(repo, orchestrator)
}

func (server *Server) initUserHandler(service *service.UserService) *handler.UserHandler {
	return handler.NewUserHandler(service)
}

func (server *Server) startGrpcServer(userHandler *handler.UserHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterAutentificationServiceServer(grpcServer, userHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
