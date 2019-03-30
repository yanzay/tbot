# tbot

**tbot** a Telegram bot server.

> Note: this is tbot v2, you can find v1 [here](https://github.com/yanzay/tbot/tree/v1.0).

## Features

- Full Telegram Bot API support
- **Zero** dependency
- Type-safe API client with functional options
- Capture messages by regexp
- Can be used with go modules
- Support for external logger
- MIT licensed

## Installation

```bash
go get github.com/yanzay/tbot
```

Go modules supported.

## Support

Join [telegram group](https://t.me/tbotgo) to get support or just to say thank you.

## Usage

Simple usage example:

[embedmd]:# (examples/basic/main.go)
```go
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
```

## Examples

Please take a look inside [examples](https://github.com/yanzay/tbot/tree/master/examples) folder.
