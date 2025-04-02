package scheduler

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/robfig/cron/v3"
	tele "gopkg.in/telebot.v4"
)

func Register(b *tele.Bot) {
	c := cron.New()
	chatID, err := strconv.Atoi(os.Getenv("CHAT_ID"))

	if err != nil {
		panic(err)
	}

	c.AddFunc("@every 1m", func() {
		b.Send(tele.ChatID(chatID), "Hello!")
	})

	slog.Info("corns registered")

	c.Start()
}
