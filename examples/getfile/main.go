package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/yanzay/tbot/v2"
)

func main() {
	token := os.Getenv("TELEGRAM_TOKEN")
	bot := tbot.New(token)
	client := bot.Client()
	bot.HandleMessage("", func(m *tbot.Message) {
		// here we check if message contains Document
		// you could also check for other types of files:
		// Audio, Photo, Video, etc.
		if m.Document != nil {
			doc, err := client.GetFile(m.Document.FileID)
			if err != nil {
				log.Println(err)
				return
			}
			url := client.FileURL(doc)
			resp, err := http.Get(url)
			if err != nil {
				log.Println(err)
				return
			}
			defer resp.Body.Close()
			out, err := os.Create(m.Document.FileName)
			if err != nil {
				log.Println(err)
				return
			}
			defer out.Close()
			io.Copy(out, resp.Body)
		}
	})
	log.Fatal(bot.Start())
}
