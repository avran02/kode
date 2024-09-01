package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server
	DB
	YandexSpeller
	JWTSecret string
}

type Server struct {
	Host     string
	Port     string
	LogLevel string
}

type DB struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type YandexSpeller struct {
	URL      string
	Language string
	Options  string
}

func New() *Config {
	if os.Getenv("ENV") != "docker" {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading .env file: %s", err.Error())
		}
	}

	return &Config{
		JWTSecret: os.Getenv("JWT_SECRET"),
		Server: Server{
			Host:     os.Getenv("SERVER_HOST"),
			Port:     os.Getenv("SERVER_PORT"),
			LogLevel: os.Getenv("LOG_LEVEL"),
		},
		DB: DB{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Database: os.Getenv("DB_DATABASE"),
		},
		YandexSpeller: YandexSpeller{
			URL:      os.Getenv("YANDEX_SPELLER_URL"),
			Language: os.Getenv("YANDEX_SPELLER_LANGUAGE"),
			Options:  os.Getenv("YANDEX_SPELLER_OPTIONS"),
		},
	}
}
