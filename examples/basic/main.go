package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/yanzay/tbot"
)

func main() {
	bot := tbot.New(os.Getenv("TELEGRAM_TOKEN"))
	c := bot.Client()
	bot.HandleMessage(".*yo.*", func(m *tbot.Message) {
		fmt.Println(m.Text)
		c.SendChatAction(m.Chat.ID, tbot.ActionTyping)
		time.Sleep(1 * time.Second)
		c.SendMessage(m.Chat.ID, "hello!")
		fmt.Println(c.GetUserProfilePhotos(m.From.ID))
	})
	log.Fatal(bot.Start())
}
