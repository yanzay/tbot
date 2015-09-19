package tbot

import (
	"fmt"
	"log"

	"github.com/tucnak/telebot"
)

type MessageType int

const (
	MessageText MessageType = iota
	MessageSticker
	MessagePhoto
	MessageAudio
	MessageVideo
	MessageLocation
	MessageDocument
)

type MessageVars map[string]string

type Message struct {
	telebot.Message
	Vars MessageVars

	replies chan *ReplyMessage
	close   chan struct{}
}

type ReplyMessage struct {
	telebot.Message
	messageType MessageType
}

func (m *Message) Reply(reply string) {
	message := &ReplyMessage{
		messageType: MessageText,
	}
	message.Text = reply
	m.replies <- message
}

func (m *Message) Replyf(reply string, values ...interface{}) {
	m.Reply(fmt.Sprintf(reply, values...))
}

func (m *Message) ReplySticker(filepath string) {
	file, err := telebot.NewFile(filepath)
	if err != nil {
		log.Println("Can't open file %s: %s", filepath, err.Error())
		return
	}
	message := &ReplyMessage{
		messageType: MessageSticker,
	}
	message.Sticker = telebot.Sticker{File: file}
	m.replies <- message
}
