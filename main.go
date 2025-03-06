package main

import (
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	pref := tele.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/hello", func(c tele.Context) error {
		return c.Send("Hello!")
	})

	b.Handle("/getInfo", func(c tele.Context) error {
		return c.Send("Hello!")
	})

	b.Handle("/products", func(c tele.Context) error {
		return c.Send("тухум, помидор, зелень, грудинка, оливка, бодринг, мол гошти")
	})

	slog.Info("Bot started successfully")
	b.Start()
}
