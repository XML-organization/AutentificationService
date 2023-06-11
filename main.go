package main

import (
	"autentification_service/startup"
	cfg "autentification_service/startup/config"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stderr)
	config := cfg.NewConfig()
	log.Println("Starting server Autentification Service...")
	server := startup.NewServer(config)
	server.Start()
}
