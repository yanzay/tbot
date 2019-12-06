# tbot - Telegram Bot Server [![GoDoc](https://godoc.org/github.com/yanzay/tbot?status.svg)](https://godoc.org/github.com/yanzay/tbot) [![Go Report Card](https://goreportcard.com/badge/github.com/yanzay/tbot)](https://goreportcard.com/report/github.com/yanzay/tbot) [![GitHub Actions](https://github.com/yanzay/tbot/workflows/Test/badge.svg)](https://github.com/yanzay/tbot/actions)

![logo](https://raw.githubusercontent.com/yanzay/tbot/master/logo.png)

## Features

- Full Telegram Bot API **4.4** support
- **Zero** dependency
- Type-safe API client with functional options
- Capture messages by regexp
- Middlewares support
- Can be used with go modules
- Support for external logger
- MIT licensed

## Installation

With go modules:

```bash
go get github.com/yanzay/tbot/v2
```

Without go modules:

```bash
go get github.com/yanzay/tbot
```

## Support

Join [telegram group](https://t.me/tbotgo) to get support or just to say thank you.

## Documentation

Documentation: [https://yanzay.github.io/tbot-doc/](https://yanzay.github.io/tbot-doc/).

Full specification: [godoc](https://godoc.org/github.com/yanzay/tbot).

## Usage

Simple usage example:

[embedmd]:# (examples/basic/main.go)
```go
package main

import (
	"log"
	"os"
	"time"

	"github.com/yanzay/tbot/v2"
)

func main() {
	bot := tbot.New(os.Getenv("TELEGRAM_TOKEN"))
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
```

## Examples

Please take a look inside [examples](https://github.com/yanzay/tbot/tree/master/examples) folder.
