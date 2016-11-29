# tbot - Telegram Bot Server [![Build Status](https://travis-ci.org/yanzay/tbot.svg?branch=master)](https://travis-ci.org/yanzay/tbot) [![Go Report Card](https://goreportcard.com/badge/github.com/yanzay/tbot)](https://goreportcard.com/report/github.com/yanzay/tbot)
[![GoDoc](https://godoc.org/github.com/yanzay/tbot?status.svg)](https://godoc.org/github.com/yanzay/tbot)

**tbot** is a Telegram bot server.

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

    err = bot.ListenAndServe()
    log.Fatal(err)
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

See full documentation here: https://godoc.org/github.com/yanzay/tbot
