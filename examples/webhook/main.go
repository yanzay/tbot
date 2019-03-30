package main

import (
	"log"
	"os"

	"github.com/yanzay/tbot"
)

func main() {
	bot := tbot.New(os.Getenv("TELEGRAM_TOKEN"),
		tbot.WithWebhook("https://test.com", ":8080"))
	c := bot.Client()
	bot.HandleMessage("ping", func(m *tbot.Message) {
		c.SendMessage(m.Chat.ID, "pong")
	})
	log.Fatal(bot.Start())
}
