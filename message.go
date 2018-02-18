package tbot

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/yanzay/tbot/model"
)

// MessageVars is a parsed message variables lookup table
type MessageVars map[string]string

// Message is a received message from chat, with parsed variables
type Message struct {
	*model.Message
	Vars         MessageVars
	replyChannel chan *model.Message
}

// MessageOption is a functional option for text messages
type MessageOption func(*model.Message)

// DisablePreview option disables web page preview when sending links.
var DisablePreview = func(msg *model.Message) {
	msg.DisablePreview = true
}

// WithMarkdown option enables Markdown style formatting for text messages.
var WithMarkdown = func(msg *model.Message) {
	msg.Markdown = true
}

// Text returns message text
func (m *Message) Text() string {
	return m.Data
}

// Reply to the user with plain text
func (m *Message) Reply(reply string, options ...MessageOption) {
	msg := &model.Message{
		ChatID: m.ChatID,
		Type:   model.MessageText,
		Data:   reply,
	}
	for _, option := range options {
		option(msg)
	}
	m.sendReply(msg)
}

// Replyf is a formatted reply to the user with plain text, with parameters like in fmt.Printf
func (m *Message) Replyf(reply string, values ...interface{}) {
	m.Reply(fmt.Sprintf(reply, values...))
}

// ReplySticker sends sticker to the chat.
func (m *Message) ReplySticker(filepath string) {
	msg := &model.Message{
		Type:   model.MessageSticker,
		Data:   filepath,
		ChatID: m.ChatID,
	}
	m.sendReply(msg)
}

// ReplyPhoto sends photo to the chat. Has optional caption.
func (m *Message) ReplyPhoto(filepath string, caption ...string) {
	msg := &model.Message{
		Type:   model.MessagePhoto,
		Data:   filepath,
		ChatID: m.ChatID,
	}
	if len(caption) > 0 {
		msg.Caption = caption[0]
	}
	m.sendReply(msg)
}

// ReplyAudio sends audio file to chat
func (m *Message) ReplyAudio(filepath string) {
	msg := &model.Message{
		Type:   model.MessageAudio,
		Data:   filepath,
		ChatID: m.ChatID,
	}
	m.sendReply(msg)
}

// ReplyDocument sends generic file (not audio, voice, image) to the chat
func (m *Message) ReplyDocument(filepath string) {
	msg := &model.Message{
		Type:   model.MessageDocument,
		Data:   filepath,
		ChatID: m.ChatID,
	}
	m.sendReply(msg)
}

// KeyboardOption is a functional option for custom keyboards
type KeyboardOption func(*model.Message)

// OneTimeKeyboard option sends keyboard that hides after the user use it once.
var OneTimeKeyboard = func(msg *model.Message) {
	msg.OneTimeKeyboard = true
}

// ReplyKeyboard sends custom reply keyboard to the user.
func (m *Message) ReplyKeyboard(text string, buttons [][]string, options ...KeyboardOption) {
	msg := &model.Message{
		Type:    model.MessageKeyboard,
		Data:    text,
		Buttons: buttons,
		ChatID:  m.ChatID,
	}
	for _, option := range options {
		option(msg)
	}
	m.sendReply(msg)
}

// ReplyLocation sends location reply to the user.
func (m *Message) ReplyLocation(longitude, latitude float64) {
	msg := &model.Message{
		Type: model.MessageLocation,
		Location: model.Location{
			Longitude: longitude,
			Latitude:  latitude,
		},
		ChatID: m.ChatID,
	}
	m.sendReply(msg)
}

// RequestContactButton sends custom reply contact button to the user.
func (m *Message) RequestContactButton(text string, button string, options ...KeyboardOption) {
	msg := &model.Message{
		Type:          model.MessageContactButton,
		Data:          text,
		ContactButton: button,
		ChatID:        m.ChatID,
	}
	for _, option := range options {
		option(msg)
	}
	m.sendReply(msg)
}

// RequestLocationButton sends custom reply location keyboard to the user.
func (m *Message) RequestLocationButton(text string, button string, options ...KeyboardOption) {
	msg := &model.Message{
		Type:           model.MessageLocationButton,
		Data:           text,
		LocationButton: button,
		ChatID:         m.ChatID,
	}
	for _, option := range options {
		option(msg)
	}
	m.sendReply(msg)
}

// InlineKeyboardButtonsOption is a functional option for inline keyboard buttons
type InlineKeyboardButtonsOption func(*model.Message)

// WithDataInlineButtons option send inline keyboard buttons with data for catch a callback.
var WithDataInlineButtons = func(msg *model.Message) {
	msg.WithDataInlineButtons = true
}

// WithURLInlineButtons option send inline keyboard buttons as url.
var WithURLInlineButtons = func(msg *model.Message) {
	msg.WithURLInlineButtons = true
}

// ReplyInlineKeyboard sends custom inline reply keyboard waiting for data callback to the user.
func (m *Message) ReplyInlineKeyboard(text string, inlineButtons []map[string]string, options ...InlineKeyboardButtonsOption) {
	msg := &model.Message{
		Type:          model.MessageInlineKeyboard,
		Data:          text,
		InlineButtons: inlineButtons,
		ChatID:        m.ChatID,
	}
	for _, option := range options {
		option(msg)
	}
	m.sendReply(msg)
}

// Download file from FileHandler
func (m *Message) Download(dir string) error {
	if m.Type != model.MessageDocument {
		return fmt.Errorf("Nothing to download")
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			return fmt.Errorf("Can't create directory for user uploads: %q", err)
		}
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

// SetReplyChannel sets channel for custom reply handling, e. g. for tests
func (m *Message) SetReplyChannel(ch chan *model.Message) {
	m.replyChannel = ch
}

func (m *Message) sendReply(msg *model.Message) {
	if m.replyChannel != nil {
		m.replyChannel <- msg
	} else {
		m.Replies <- msg
	}
}
