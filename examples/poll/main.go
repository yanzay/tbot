package main

import (
	"fmt"
	"os"

	"github.com/yanzay/tbot"
)

var client *tbot.Client

func main() {
	bot := tbot.New(os.Getenv("TELEGRAM_TOKEN"))
	client = bot.Client()
	// listen poll message and send poll
	bot.HandleMessage("poll", sendPoll)
	// handle poll updates, just print on the screen
	bot.HandlePollUpdate(func(p *tbot.Poll) {
		fmt.Println("Poll update received:")
		fmt.Println(p.Question)
		for _, opt := range p.Options {
			fmt.Println(opt.Text, opt.VoterCount)
		}
	})
	bot.Start()
}

func sendPoll(m *tbot.Message) {
	options := []string{
		"Perfect",
		"Good",
		"So so",
	}
	client.SendPoll(m.Chat.ID, "How are you?", options)
}
