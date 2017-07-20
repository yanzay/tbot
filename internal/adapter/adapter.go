package adapter

import (
	"fmt"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/yanzay/tbot/model"
)

type BotAdapter interface {
	Send(*model.Message) error
	GetUpdatesChan(string, string) (<-chan *model.Message, error)
	GetUserName() string
	GetFirstName() string
}

type Bot struct {
	tbot *tgbotapi.BotAPI
}

func CreateBot(token string) (BotAdapter, error) {
	tbot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &Bot{tbot: tbot}, nil
}

func (b *Bot) Send(m *model.Message) error {
	c := chattableFromMessage(m)
	if c != nil {
		_, err := b.tbot.Send(c)
		return err
	}
	return fmt.Errorf("Trying to send nil chattable. Message: %v", m)
}

func (b *Bot) GetUpdatesChan(webhookURL string, listenAddr string) (<-chan *model.Message, error) {
	messages := make(chan *model.Message)
	var updates <-chan tgbotapi.Update
	var err error
	if webhookURL == "" {
		b.tbot.RemoveWebhook()
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		updates, err = b.tbot.GetUpdatesChan(u)
		if err != nil {
			return nil, err
		}
	} else {
		config := tgbotapi.NewWebhook(webhookURL)
		b.tbot.SetWebhook(config)
		updates = b.tbot.ListenForWebhook("/")
		go http.ListenAndServe(listenAddr, nil)
	}
	go b.adaptUpdates(updates, messages)
	return messages, nil
}

func (b *Bot) GetUserName() string {
	return b.tbot.Self.UserName
}

func (b *Bot) GetFirstName() string {
	return b.tbot.Self.FirstName
}

func (b *Bot) adaptUpdates(updates <-chan tgbotapi.Update, messages chan<- *model.Message) {
	var err error
	for update := range updates {
		var updateMessage *tgbotapi.Message

		if update.Message != nil {
			updateMessage = update.Message
		}
		if update.ChannelPost != nil {
			updateMessage = update.ChannelPost
		}
		if updateMessage == nil {
			continue
		}

		message := &model.Message{
			Replies: make(chan *model.Message),
			ChatID:  updateMessage.Chat.ID,
		}
		if updateMessage.From != nil {
			message.From = model.User{
				updateMessage.From.ID,
				updateMessage.From.FirstName,
				updateMessage.From.LastName,
				updateMessage.From.UserName,
				updateMessage.From.LanguageCode,
			}
		}
		switch {
		case updateMessage.Document != nil:
			message.Type = model.MessageDocument
			message.Data, err = b.tbot.GetFileDirectURL(updateMessage.Document.FileID)
			if err != nil {
				log.Println(err)
				continue
			}
			messages <- message
		case updateMessage.Text != "":
			message.Type = model.MessageText
			message.Data = updateMessage.Text
			messages <- message
		}
	}
}

func chattableFromMessage(m *model.Message) tgbotapi.Chattable {
	switch m.Type {
	case model.MessageText:
		msg := tgbotapi.NewMessage(m.ChatID, m.Data)
		msg.DisableWebPagePreview = m.DisablePreview
		if m.Markdown {
			msg.ParseMode = tgbotapi.ModeMarkdown
		}
		return msg
	case model.MessageSticker:
		return tgbotapi.NewStickerUpload(m.ChatID, m.Data)
	case model.MessagePhoto:
		photo := tgbotapi.NewPhotoUpload(m.ChatID, m.Data)
		photo.Caption = m.Caption
		return photo
	case model.MessageAudio:
		return tgbotapi.NewAudioUpload(m.ChatID, m.Data)
	case model.MessageDocument:
		return tgbotapi.NewDocumentUpload(m.ChatID, m.Data)
	case model.MessageKeyboard:
		msg := tgbotapi.NewMessage(m.ChatID, m.Data)
		btns := buttonsFromStrings(m.Buttons)
		keyboard := tgbotapi.NewReplyKeyboard(btns...)
		keyboard.OneTimeKeyboard = m.OneTimeKeyboard
		msg.ReplyMarkup = keyboard
		if m.Markdown {
			msg.ParseMode = tgbotapi.ModeMarkdown
		}
		return msg
	}
	return nil
}

func buttonsFromStrings(strs [][]string) [][]tgbotapi.KeyboardButton {
	btns := make([][]tgbotapi.KeyboardButton, len(strs))
	for i, buttonRow := range strs {
		btns[i] = make([]tgbotapi.KeyboardButton, len(buttonRow))
		for j, buttonText := range buttonRow {
			btns[i][j] = tgbotapi.NewKeyboardButton(buttonText)
		}
	}
	return btns
}
