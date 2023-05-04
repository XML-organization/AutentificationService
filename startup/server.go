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

func (server *Server) Start() {
	postgresClient := server.initPostgresClient()
	userRepo := server.initUserRepository(postgresClient)
	userService := server.initUserService(userRepo)
	userHandler := server.initUserHandler(userService)

	server.startGrpcServer(userHandler)
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

func (server *Server) initUserService(repo *repository.UserRepository) *service.UserService {
	return service.NewUserService(repo)
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
