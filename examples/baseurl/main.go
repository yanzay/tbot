package main

import (
	"log"
	"os"
	"time"

	"github.com/yanzay/tbot/v2"
)

/* Before setup nginx as proxy. See https://habr.com/ru/post/424427/

server {
  listen 443 ssl http2;
  server_name my-telegram-proxy.server;

  # SSL options skipped

  location / {
    proxy_set_header X-Forwarded-Host $host;
    proxy_set_header X-Forwarded-Server $host;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_pass https://api.telegram.org/;
    client_max_body_size 100M;
  }
}
*/

func main() {
	bot := tbot.New(os.Getenv("TELEGRAM_TOKEN"),
		tbot.WithBaseURL("https://my-telegram-proxy.server"))
	c := bot.Client()
	bot.HandleMessage(".*yo.*", func(m *tbot.Message) {
		c.SendChatAction(m.Chat.ID, tbot.ActionTyping)
		time.Sleep(1 * time.Second)
		c.SendMessage(m.Chat.ID, "hello!")
	})
	err := bot.Start()
	if err != nil {
		log.Fatal(err)
	}
}
