package tbot

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/yanzay/tbot/internal/adapter"
)

// MessageVars is a parsed message variables lookup table
type MessageVars map[string]string

// Message is a received message from chat, with parsed variables
type Message struct {
	*adapter.Message
	Vars MessageVars
}

type Option func(*adapter.Message)

var DisablePreview = func(msg *adapter.Message) {
	msg.DisablePreview = true
}

// Text returns message text
func (m *Message) Text() string {
	return m.Data
}

// Reply to the user with plain text
func (m *Message) Reply(reply string, options ...Option) {
	msg := &adapter.Message{
		ChatID: m.ChatID,
		Type:   adapter.MessageText,
		Data:   reply,
	}
	for _, option := range options {
		option(msg)
	}
	m.Replies <- msg
}

// Replyf is a formatted reply to the user with plain text, with parameters like in fmt.Printf
func (m *Message) Replyf(reply string, values ...interface{}) {
	m.Reply(fmt.Sprintf(reply, values...))
}

// ReplySticker sends sticker to the chat.
func (m *Message) ReplySticker(filepath string) {
	msg := &adapter.Message{
		Type:   adapter.MessageSticker,
		Data:   filepath,
		ChatID: m.ChatID,
	}
	m.Replies <- msg
}

// ReplyPhoto sends photo to the chat. Has optional caption.
func (m *Message) ReplyPhoto(filepath string, caption ...string) {
	msg := &adapter.Message{
		Type:   adapter.MessagePhoto,
		Data:   filepath,
		ChatID: m.ChatID,
	}
	if len(caption) > 0 {
		msg.Caption = caption[0]
	}
	m.Replies <- msg
}

// ReplyAudio sends audio file to chat
func (m *Message) ReplyAudio(filepath string) {
	msg := &adapter.Message{
		Type:   adapter.MessageAudio,
		Data:   filepath,
		ChatID: m.ChatID,
	}
	m.Replies <- msg
}

// ReplyDocument sends generic file (not audio, voice, image) to the chat
func (m *Message) ReplyDocument(filepath string) {
	msg := &adapter.Message{
		Type:   adapter.MessageDocument,
		Data:   filepath,
		ChatID: m.ChatID,
	}
	m.Replies <- msg
}

// Download file from FileHandler
func (m *Message) Download(dir string) error {
	if m.Type != adapter.MessageDocument {
		return fmt.Errorf("Nothing to download")
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0755)
	}

	tokens := strings.Split(m.Vars["url"], "/")
	fileName := tokens[len(tokens)-1]

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
