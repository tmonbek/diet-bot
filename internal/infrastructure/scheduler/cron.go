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

	c.AddFunc("@hourly", func() {
		b.Send(tele.ChatID(chatID), "Suv ichish vaqti boldi")
	})

	c.AddFunc("30 7 * * *", func() {
		b.Send(tele.ChatID(chatID), "Assalamu alekum")
	})

	c.AddFunc("0 8 * * *", func() {
		b.Send(tele.ChatID(chatID), "Nonushta vaqti boldi: 4 ta tuhum oqi 2 ta sarigi blan va 50g ovsanka")
	})

	c.AddFunc("0 13 * * *", func() {
		b.Send(tele.ChatID(chatID), "Tushlik vaqti boldi: 200g tovuq va sabzavotli salat")
	})

	c.AddFunc("0 19 * * *", func() {
		b.Send(tele.ChatID(chatID), "Kechki ovqat vaqti boldi: 200g mol goshti  yoki baliq va salat bodring pomidor")
	})
	
	c.AddFunc("30 22 * * *", func() {
		b.Send(tele.ChatID(chatID), "кеч 22:30 ларда! 2 та кайнатилган тухум оки  1 порц вей изолят ичига 5 гр аргенин солиб  250 мл сувга аралаштириб ичамиз")
	})

	c.AddFunc("0 16 * * *", func() {
		b.Send(tele.ChatID(chatID), "Perekus protein vey izolyat ichamiz")
	})

	slog.Info("corns registered")

	c.Start()
}
