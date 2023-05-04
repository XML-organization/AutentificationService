package config

import "os"

type Config struct {
	Port                  string
	AutentificationDBHost string
	AutentificationDBPort string
	AutentificationDBName string
	AutentificationDBUser string
	AutentificationDBPass string
}

func NewConfig() *Config {
	return &Config{
		Port:                  os.Getenv("AUTENTIFICATION_SERVICE_PORT"),
		AutentificationDBHost: os.Getenv("AUTENTIFICATION_DB_HOST"),
		AutentificationDBPort: os.Getenv("AUTENTIFICATION_DB_PORT"),
		AutentificationDBName: os.Getenv("AUTENTIFICATION_DB_NAME"),
		AutentificationDBUser: os.Getenv("AUTENTIFICATION_DB_USER"),
		AutentificationDBPass: os.Getenv("AUTENTIFICATION_DB_PASS"),
	}
}
