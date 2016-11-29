package tbot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type MessageVars map[string]string

// Message is a received message from chat, with parsed variables
type Message struct {
	tgbotapi.Message
	Vars MessageVars

	replies chan *ReplyMessage
	close   chan struct{}
}

// ReplyMessage is a bot message
type ReplyMessage struct {
	msg tgbotapi.Chattable
}

// Reply to the user with plain text
func (m Message) Reply(reply string) {
	message := &ReplyMessage{
		msg: tgbotapi.NewMessage(m.Chat.ID, reply),
	}
	m.replies <- message
}

// Replyf is a formatted reply to the user with plain text, with parameters like in Printf
func (m Message) Replyf(reply string, values ...interface{}) {
	m.Reply(fmt.Sprintf(reply, values...))
}

func (m Message) ReplySticker(filepath string) {
	message := &ReplyMessage{}
	msg := tgbotapi.NewStickerUpload(m.Chat.ID, filepath)
	message.msg = msg
	m.replies <- message
}

// ReplyPhoto sends photo to the chat. Has optional caption.
func (m Message) ReplyPhoto(filepath string, caption ...string) {
	message := &ReplyMessage{}
	msg := tgbotapi.NewPhotoUpload(m.Chat.ID, filepath)
	if len(caption) > 0 {
		msg.Caption = caption[0]
	}
	message.msg = msg
	m.replies <- message
}

// ReplyAudio sends audio file to chat
func (m Message) ReplyAudio(filepath string) {
	message := &ReplyMessage{}
	msg := tgbotapi.NewAudioUpload(m.Chat.ID, filepath)
	message.msg = msg
	m.replies <- message
}

// ReplyDocument sends generic file (not audio, voice, image) to the chat
func (m Message) ReplyDocument(filepath string) {
	message := &ReplyMessage{}
	msg := tgbotapi.NewDocumentUpload(m.Chat.ID, filepath)
	message.msg = msg
	m.replies <- message
}
