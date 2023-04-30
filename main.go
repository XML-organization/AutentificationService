package main

import (
	"autentification_service/handler"
	"autentification_service/model"
	"autentification_service/repository"
	"autentification_service/service"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	connectionStr := "host=localhost user=postgres password=password dbname=AuthenticationDatabase port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(connectionStr), &gorm.Config{})
	if err != nil {
		print(err)
		return nil
	}

	database.AutoMigrate(&model.UserCredentials{})

	return database
}

func GetClient(host, user, password, dbname, port string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func initPostgresClient() *gorm.DB {
	client, err := GetClient(
		os.Getenv("AUTENTIFICATION_DB_HOST"), os.Getenv("AUTENTIFICATION_DB_USER"),
		os.Getenv("AUTENTIFICATION_DB_PASS"), os.Getenv("AUTENTIFICATION_DB_NAME"),
		os.Getenv("AUTENTIFICATION_DB_PORT"))
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func startServer(userHandler *handler.UserHandler) {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/create", userHandler.Create).Methods("POST")
	router.HandleFunc("/login", userHandler.Login).Methods("POST")
	router.HandleFunc("/getLoggedUser", userHandler.User).Methods("GET")
	router.HandleFunc("/logout", userHandler.Logout).Methods("POST")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://192.168.0.17:5173", "http://localhost:5173", "http://192.168.137.1:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	handler1 := corsHandler.Handler(router)

	router.Methods("OPTIONS").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.WriteHeader(http.StatusNoContent)
		})

	println("Server starting")
	log.Fatal(http.ListenAndServe(":8082", handler1))
}

func main() {
	//database := initDB()
	database := initPostgresClient()
	if database == nil {
		log.Fatal("FAILED TO CONNECT TO DB")
	}
	repoUser := &repository.UserRepository{DatabaseConnection: database}
	serviceUser := &service.UserService{UserRepo: repoUser}
	handlerUser := &handler.UserHandler{UserService: serviceUser}

	startServer(handlerUser)
}
