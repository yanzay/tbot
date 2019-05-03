package main

import (
	"log"
	"os"
	"time"

	"github.com/yanzay/tbot/v2"
)

func stat(h tbot.UpdateHandler) tbot.UpdateHandler {
	return func(u *tbot.Update) {
		start := time.Now()
		h(u)
		log.Printf("Handle time: %v", time.Now().Sub(start))
	}
}

func main() {
	bot := tbot.New(os.Getenv("TELEGRAM_TOKEN"))
	c := bot.Client()
	bot.Use(stat) // add stat middleware to bot
	bot.HandleMessage("", func(m *tbot.Message) {
		c.SendMessage(m.Chat.ID, "hello!")
	})
	err := bot.Start()
	if err != nil {
		log.Fatal(err)
	}
}
