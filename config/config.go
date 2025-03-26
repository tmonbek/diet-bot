package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found")
	}
}

func Settings() tele.Settings {
	return tele.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
}
