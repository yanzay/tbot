package main

import (
	"fmt"
	"os"

	"github.com/yanzay/tbot/v2"
)

type application struct {
	client  *tbot.Client
	votings map[string]*voting
}

type voting struct {
	ups   int
	downs int
}

func main() {
	token := os.Getenv("TELEGRAM_TOKEN")
	bot := tbot.New(token)
	app := &application{
		votings: make(map[string]*voting),
	}
	app.client = bot.Client()
	bot.HandleMessage("/vote", app.votingHandler)
	bot.HandleCallback(app.callbackHandler)
	bot.Start()
}

func (a *application) votingHandler(m *tbot.Message) {
	buttons := makeButtons(0, 0)
	msg, _ := a.client.SendMessage(m.Chat.ID, "Please vote", tbot.OptInlineKeyboardMarkup(buttons))
	votingID := fmt.Sprintf("%s:%d", m.Chat.ID, msg.MessageID)
	a.votings[votingID] = &voting{}
}

func (a *application) callbackHandler(cq *tbot.CallbackQuery) {
	votingID := fmt.Sprintf("%s:%d", cq.Message.Chat.ID, cq.Message.MessageID)
	v := a.votings[votingID]
	if cq.Data == "up" {
		v.ups++
	}
	if cq.Data == "down" {
		v.downs++
	}
	buttons := makeButtons(v.ups, v.downs)
	a.client.EditMessageReplyMarkup(cq.Message.Chat.ID, cq.Message.MessageID, tbot.OptInlineKeyboardMarkup(buttons))
	a.client.AnswerCallbackQuery(cq.ID, tbot.OptText("OK"))
}

func makeButtons(ups, downs int) *tbot.InlineKeyboardMarkup {
	km := tbot.NewKeyboardMaker()
	firstRow := km.AddRow()
	firstRow.AddButton(fmt.Sprintf("👍 %d", ups), "up", "")
	firstRow.AddButton(fmt.Sprintf("👎 %d", downs), "down", "")

	secondRow := km.AddRow()
	secondRow.AddButton("To a random web!", "rand-web", "http://google.com")

	return km.Build()
}
