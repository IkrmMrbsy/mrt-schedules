package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	MRTApiURL  string
	ServerPort string
}

func LoadConfig() *config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, fallback to system env")
	}

	return &config{
		MRTApiURL:  os.Getenv("MRT_API_URL"),
		ServerPort: os.Getenv("SERVER_PORT"),
	}
}
