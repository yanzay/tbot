package main

import (
	"fmt"
	"log"
	"os"

	"github.com/yanzay/tbot/v2"
)

func main() {
	bot := tbot.New(os.Getenv("TELEGRAM_TOKEN"))
	c := bot.Client()
	bot.HandleMessage("image", func(m *tbot.Message) {
		_, err := c.SendPhotoFile(m.Chat.ID, "image.png", tbot.OptCaption("this is image"))
		if err != nil {
			fmt.Println(err)
		}
	})
	err := bot.Start()
	if err != nil {
		log.Fatal(err)
	}
}
