package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Telegram struct {
		Token   string
		Timeout time.Duration
	}
}

func LoadConfig() Config {
	err := godotenv.Load()
	var cfg Config

	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}

	cfg.Telegram.Token = os.Getenv("TELEGRAM_TOKEN")
	cfg.Telegram.Timeout = 10 * time.Second
	return cfg
}
