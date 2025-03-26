package handlers

import (
	"log/slog"

	tele "gopkg.in/telebot.v4"
)

func Register(b *tele.Bot) {
	b.Handle("/hello", func(c tele.Context) error {
		return c.Send("Hello!")
	})

	b.Handle("/getInfo", func(c tele.Context) error {
		return c.Send("Hello!")
	})

	b.Handle("/products", func(c tele.Context) error {
		return c.Send("тухум, помидор, зелень, грудинка, оливка, бодринг, мол гошти")
	})

	slog.Info("handlers registered")
}
