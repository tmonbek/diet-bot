package main

import (
	"diet-bot/config"
	telegram "diet-bot/internal/application"
	"diet-bot/internal/handlers"
)

func main() {
	config.LoadEnv()
	pref := config.Settings()

	b := telegram.NewBot(pref)

	handlers.Register(b)

	b.Start()
}
