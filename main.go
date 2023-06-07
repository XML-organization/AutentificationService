package main

import (
	"autentification_service/startup"
	cfg "autentification_service/startup/config"
)

func main() {
	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
