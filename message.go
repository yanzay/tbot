package tbot

import (
	"fmt"

	"github.com/tucnak/telebot"
)

type MessageVars map[string]string

type Message struct {
	telebot.Message
	Vars MessageVars

	replies chan string
	close   chan struct{}
}

func (m *Message) Reply(reply string) {
	m.replies <- reply
}

func (m *Message) Replyf(reply string, values ...interface{}) {
	m.Reply(fmt.Sprintf(reply, values...))
}
