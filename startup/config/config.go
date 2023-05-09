package config

import "os"

type Config struct {
	Port                         string
	AutentificationDBHost        string
	AutentificationDBPort        string
	AutentificationDBName        string
	AutentificationDBUser        string
	AutentificationDBPass        string
	NatsHost                     string
	NatsPort                     string
	NatsUser                     string
	NatsPass                     string
	CreateUserCommandSubject     string
	CreateUserReplySubject       string
	ChangePasswordCommandSubject string
	ChangePasswordReplySubject   string
	UpdateUserCommandSubject     string
	UpdateUserReplySubject       string
}

func NewConfig() *Config {
	return &Config{
		Port:                         os.Getenv("AUTENTIFICATION_SERVICE_PORT"),
		AutentificationDBHost:        os.Getenv("AUTENTIFICATION_DB_HOST"),
		AutentificationDBPort:        os.Getenv("AUTENTIFICATION_DB_PORT"),
		AutentificationDBName:        os.Getenv("AUTENTIFICATION_DB_NAME"),
		AutentificationDBUser:        os.Getenv("AUTENTIFICATION_DB_USER"),
		AutentificationDBPass:        os.Getenv("AUTENTIFICATION_DB_PASS"),
		NatsHost:                     os.Getenv("NATS_HOST"),
		NatsPort:                     os.Getenv("NATS_PORT"),
		NatsUser:                     os.Getenv("NATS_USER"),
		NatsPass:                     os.Getenv("NATS_PASS"),
		CreateUserCommandSubject:     os.Getenv("CREATE_USER_COMMAND_SUBJECT"),
		CreateUserReplySubject:       os.Getenv("CREATE_USER_REPLY_SUBJECT"),
		ChangePasswordCommandSubject: os.Getenv("CHANGE_PASSWORD_COMMAND_SUBJECT"),
		ChangePasswordReplySubject:   os.Getenv("CHANGE_PASSWORD_REPLY_SUBJECT"),
		UpdateUserCommandSubject:     os.Getenv("UPDATE_USER_COMMAND_SUBJECT"),
		UpdateUserReplySubject:       os.Getenv("UPDATE_USER_REPLY_SUBJECT"),
	}
}
