package main

import (
	"os"
	"time"

	"github.com/yanzay/tbot"
)

func main() {
	bot := tbot.New(os.Getenv("TELEGRAM_TOKEN"))
	c := bot.Client()
	bot.HandleMessage(".*yo.*", func(m *tbot.Message) {
		c.SendChatAction(m.Chat.ID, tbot.ActionTyping)
		time.Sleep(1 * time.Second)
		c.SendMessage(m.Chat.ID, "hello!")
	})
	bot.Start()
}
