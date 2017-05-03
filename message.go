package tbot

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// MessageVars is a parsed message variables lookup table
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

type Option func(*tgbotapi.MessageConfig)

var DisablePreview = func(msg *tgbotapi.MessageConfig) {
	msg.DisableWebPagePreview = true
}

// Reply to the user with plain text
func (m Message) Reply(reply string, options ...Option) {
	msg := tgbotapi.NewMessage(m.Chat.ID, reply)
	for _, option := range options {
		option(&msg)
	}
	message := &ReplyMessage{msg: msg}
	m.replies <- message
}

// Replyf is a formatted reply to the user with plain text, with parameters like in fmt.Printf
func (m Message) Replyf(reply string, values ...interface{}) {
	m.Reply(fmt.Sprintf(reply, values...))
}

// ReplySticker sends sticker to the chat.
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

// Download file from FileHandler
func (m Message) Download(dir string) error {
	if m.Document == nil {
		return fmt.Errorf("Nothing to download")
	}

	fileName := m.Document.FileName
	if fileName == "" {
		tokens := strings.Split(m.Vars["url"], "/")
		fileName = tokens[len(tokens)-1]
	}

	file, err := os.Create(filepath.Join(dir, fileName))
	if err != nil {
		return fmt.Errorf("Error creating file: %q", err)
	}
	defer file.Close()

	resp, err := http.Get(m.Vars["url"])
	if err != nil {
		return fmt.Errorf("Error downloading from %s: %q", m.Vars["url"], err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("[Tbot] Error downloading file: %q", err)
	}
	return nil
}
