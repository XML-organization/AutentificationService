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

	commandPublisher := server.initPublisher(server.config.CreateUserCommandSubject)
	replySubscriber := server.initSubscriber(server.config.CreateUserReplySubject, QueueGroup)
	createUserOrchestrator := server.initCreateOrderOrchestrator(commandPublisher, replySubscriber)

	userService := server.initUserService(userRepo, createUserOrchestrator)
	userHandler := server.initUserHandler(userService)

	server.startGrpcServer(userHandler)
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
