package main

import (
	"diet-bot/internal/domain/handlers"
	"diet-bot/internal/infrastructure/config"
	"diet-bot/internal/infrastructure/scheduler"
	"diet-bot/internal/infrastructure/telegram"
	postgres "diet-bot/internal/store"
)

func main() {
	config.LoadEnv()

	pref := config.Settings()

	postgres.InitDB()
	defer postgres.CloseDB()

	b := telegram.NewBot(pref)

	handlers.Register(b)

	scheduler.Register(b)

	b.Start()
}
