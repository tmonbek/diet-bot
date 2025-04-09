package scheduler

import (
	postgres "diet-bot/internal/store"
	"log"
	"log/slog"
	"time"

	"github.com/robfig/cron/v3"
	"gopkg.in/telebot.v4"
	tele "gopkg.in/telebot.v4"
)

func Register(b *tele.Bot) {
	c := cron.New()
	var chats = postgres.GetAllChats()

	c.AddFunc("@every 2h", func() {
		now := time.Now()
		hour := now.Hour()

		if hour >= 7 && hour < 22 {
			broadcast(b, "Suv ichish vaqti boldi", chats)
		}
	})

	c.AddFunc("30 7 * * *", func() {
		broadcast(b, "Assalamu alekum", chats)
	})

	c.AddFunc("0 8 * * *", func() {
		broadcast(b, "Nonushta vaqti boldi: 4 ta tuhum oqi 2 ta sarigi blan va 50g ovsanka", chats)
	})

	c.AddFunc("0 13 * * *", func() {
		broadcast(b, "Tushlik vaqti boldi: 200g tovuq va sabzavotli salat", chats)
	})

	c.AddFunc("0 19 * * *", func() {
		broadcast(b, "Kechki ovqat vaqti boldi: 200g mol goshti  yoki baliq va salat bodring pomidor", chats)
	})

	c.AddFunc("30 22 * * *", func() {
		broadcast(b, "кеч 22:30 ларда! 2 та кайнатилган тухум оки  1 порц вей изолят ичига 5 гр аргенин солиб  250 мл сувга аралаштириб ичамиз", chats)
	})

	c.AddFunc("0 16 * * *", func() {
		broadcast(b, "Perekus protein vey izolyat ichamiz", chats)
	})

	slog.Info("corns registered")

	c.Start()
}

func broadcast(bot *telebot.Bot, msg string, chats []int64) {
	const workerCount = 100
	jobs := make(chan int64, len(chats))

	for i := 0; i < workerCount; i++ {
		go func() {
			for id := range jobs {
				_, err := bot.Send(telebot.ChatID(id), msg)
				if err != nil {
					log.Printf("Failed to send to %d: %v", id, err)
				}
				time.Sleep(30 * time.Millisecond)
			}
		}()
	}

	for _, id := range chats {
		jobs <- id
	}
	close(jobs)
}
