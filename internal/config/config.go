package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type config struct {
	ServerPort  string
	HttpTimeout time.Duration
	MRTApiURL   string
}

func LoadConfig() *config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, fallback to system env")
	}

	timeout, _ := strconv.Atoi(os.Getenv("HTTP_TIMEOUT"))
	if timeout == 0 {
		timeout = 10
	}

	return &config{
		ServerPort:  os.Getenv("SERVER_PORT"),
		HttpTimeout: time.Duration(timeout) * time.Second,
		MRTApiURL:   os.Getenv("MRT_API_URL"),
	}
}
