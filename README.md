# tbot - Telegram Bot Server

**tbot** is a bot server built on top of http://github.com/tucnak/telebot

It feels like net/http Server, so it's easy to use:

```go
package main

import (
    "log"
    "os"
    "time"

    "github.com/yanzay/tbot"
)

func main() {
    token := os.Getenv("TELEGRAM_TOKEN")
    bot, err := tbot.NewServer(token)
    if err != nil {
        log.Fatal(err)
    }

    bot.Handle("yo", "YO!")

    bot.HandleFunc("/hi", HiHandler)
    bot.HandleFunc("/say {text}", SayHandler)
    bot.HandleFunc("/sticker", StickerHandler)
    bot.HandleFunc("/photo", PhotoHandler)

    bot.HandleDefault(EchoHandler)

    bot.ListenAndServe()
}

func HiHandler(message tbot.Message) {
    message.Replyf("Hello, %s!", message.Sender.FirstName)
    time.Sleep(1 * time.Second)
    message.Reply("What's up?")
}

func SayHandler(message tbot.Message) {
    message.Reply(message.Vars["text"])
}

func EchoHandler(message tbot.Message) {
    message.Reply(message.Text)
}

func StickerHandler(message tbot.Message) {
    message.ReplySticker("sticker.png")
}

func PhotoHandler(message tbot.Message) {
    message.ReplyPhoto("photo.jpg", "it's me")
}
```

Now it supports only text messages, stickers and photos.
