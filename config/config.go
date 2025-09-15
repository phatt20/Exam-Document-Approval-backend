package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		App      App
		Postgres *PostgresConfig
	}

	App struct {
		Name  string
		Url   string
		Stage string
	}

	PostgresConfig struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
		SSLMode  string
		Schema   string
	}
)

func LoadConfig(path string) Config {
	if err := godotenv.Load(path); err != nil {
		log.Println("⚠️ Warning: .env not found, using system ENV only")
	}

	var pg *PostgresConfig
	if os.Getenv("PG_HOST") != "" {
		port, _ := strconv.Atoi(os.Getenv("PG_PORT"))
		pg = &PostgresConfig{
			Host:     os.Getenv("PG_HOST"),
			Port:     port,
			User:     os.Getenv("PG_USER"),
			Password: os.Getenv("PG_PASSWORD"),
			DBName:   os.Getenv("PG_DBNAME"),
			SSLMode:  os.Getenv("PG_SSLMODE"),
			Schema:   os.Getenv("PG_SCHEMA"),
		}
	}

	return Config{
		App: App{
			Name:  os.Getenv("APP_NAME"),
			Url:   os.Getenv("APP_URL"),
			Stage: os.Getenv("APP_STAGE"),
		},
		Postgres: pg,
	}
}
