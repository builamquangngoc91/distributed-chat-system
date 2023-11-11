package configs

import "os"

type Database struct {
	Host     string
	Username string
	Password string
	Name     string
	Port     string
}

type Redis struct {
	Host string
	Port string
}

type AuthService struct {
	BaseURL string
}

type ChatService struct {
	Port string
}

type Config struct {
	Database    Database
	Redis       Redis
	AuthService AuthService
	ChatService ChatService
}

var Cfg Config

func LoadConfig() {
	Cfg = Config{
		Database: Database{
			Host:     os.Getenv("SN_DB_HOST"),
			Username: os.Getenv("SN_DB_USERNAME"),
			Password: os.Getenv("SN_DB_PASSWORD"),
			Name:     os.Getenv("SN_DB_NAME"),
			Port:     os.Getenv("SN_DB_PORT"),
		},
		Redis: Redis{
			Host: os.Getenv("SN_RD_HOST"),
			Port: os.Getenv("SN_RD_PORT"),
		},
		AuthService: AuthService{
			BaseURL: os.Getenv("SN_AUTH_SERVICE_BASE_URL"),
		},
		ChatService: ChatService{
			Port: os.Getenv("SN_CHAT_SERVICE_PORT"),
		},
	}
}
