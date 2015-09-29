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
	photo       *telebot.Photo
	audio       *telebot.Audio
	document    *telebot.Document
}

func (m Message) Reply(reply string) {
	message := &ReplyMessage{
		messageType: MessageText,
	}
	message.Text = reply
	m.replies <- message
}

func (m Message) Replyf(reply string, values ...interface{}) {
	m.Reply(fmt.Sprintf(reply, values...))
}

func (m Message) ReplySticker(filepath string) {
	file, err := telebot.NewFile(filepath)
	if err != nil {
		log.Printf("Can't open file %s: %s", filepath, err.Error())
		return
	}
	message := &ReplyMessage{
		messageType: MessageSticker,
	}
	message.Sticker = telebot.Sticker{File: file}
	m.replies <- message
}

func (m Message) ReplyPhoto(filepath string, caption ...string) {
	file, err := telebot.NewFile(filepath)
	if err != nil {
		log.Printf("Can't open file %s: %s", filepath, err.Error())
		return
	}
	thumb := telebot.Thumbnail{File: file}
	message := &ReplyMessage{messageType: MessagePhoto}
	message.photo = &telebot.Photo{Thumbnail: thumb}
	if len(caption) > 0 {
		message.photo.Caption = caption[0]
	}
	m.replies <- message
}

func (m Message) ReplyAudio(filepath string) {
	file, err := telebot.NewFile(filepath)
	if err != nil {
		log.Printf("Can't open file %s: %s", filepath, err.Error())
		return
	}
	audio := telebot.Audio{File: file}
	message := &ReplyMessage{messageType: MessageAudio, audio: &audio}
	m.replies <- message
}

func (m Message) ReplyDocument(filepath string) {
	file, err := telebot.NewFile(filepath)
	if err != nil {
		log.Printf("Can't open file %s: %s", filepath, err.Error())
		return
	}
	doc := telebot.Document{File: file}
	message := &ReplyMessage{messageType: MessageDocument, document: &doc}
	m.replies <- message
}
