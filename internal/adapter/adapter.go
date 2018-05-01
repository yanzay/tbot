package adapter

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/yanzay/tbot/model"
)

type BotAdapter interface {
	Send(*model.Message) error
	SendRaw(string, map[string]string) error
	GetUpdatesChan(string, string) (<-chan *model.Message, error)
	GetUserName() string
	GetFirstName() string
}

type Bot struct {
	tbot *tgbotapi.BotAPI
}

func CreateBot(token string, httpClient *http.Client) (BotAdapter, error) {
	tbot, err := tgbotapi.NewBotAPIWithClient(token, httpClient)
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

func (b *Bot) SendRaw(endpoint string, params map[string]string) error {
	values := url.Values{}
	for k, v := range params {
		values.Add(k, v)
	}
	_, err := b.tbot.MakeRequest(endpoint, values)
	return err
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
		if update.CallbackQuery != nil && update.CallbackQuery.Message != nil {
			updateMessage = update.CallbackQuery.Message
		}
		if updateMessage == nil {
			continue
		}

		message := &model.Message{
			Replies:     make(chan *model.Message),
			ChatID:      updateMessage.Chat.ID,
			ChatType:    updateMessage.Chat.Type,
			ForwardDate: updateMessage.ForwardDate,
		}
		if updateMessage.From != nil {
			message.From = model.User{
				ID:           updateMessage.From.ID,
				FirstName:    updateMessage.From.FirstName,
				LastName:     updateMessage.From.LastName,
				UserName:     updateMessage.From.UserName,
				LanguageCode: updateMessage.From.LanguageCode,
			}
		}
		if updateMessage.Contact != nil {
			message.Contact = model.Contact{
				PhoneNumber: updateMessage.Contact.PhoneNumber,
				FirstName:   updateMessage.Contact.FirstName,
				LastName:    updateMessage.Contact.LastName,
				UserID:      updateMessage.Contact.UserID,
			}
		}
		if updateMessage.Location != nil {
			message.Location = model.Location{
				Longitude: updateMessage.Location.Longitude,
				Latitude:  updateMessage.Location.Latitude,
			}
		}
		if update.CallbackQuery != nil {
			message.CallbackQuery = model.CallbackQuery{
				ID:              update.CallbackQuery.ID,
				InlineMessageID: update.CallbackQuery.InlineMessageID,
				ChatInstance:    update.CallbackQuery.ChatInstance,
				Data:            update.CallbackQuery.Data,
				GameShortName:   update.CallbackQuery.GameShortName,
			}
			if update.CallbackQuery.From != nil {
				message.CallbackQuery.From = model.User{
					ID:           update.CallbackQuery.From.ID,
					FirstName:    update.CallbackQuery.From.FirstName,
					LastName:     update.CallbackQuery.From.LastName,
					UserName:     update.CallbackQuery.From.UserName,
					LanguageCode: update.CallbackQuery.From.LanguageCode,
				}
			}
		}
		switch {
		case update.CallbackQuery != nil:
			message.Type = model.MessageInlineKeyboard
			if updateMessage.Text != "" {
				message.Data = updateMessage.Text
			}
			messages <- message
		case updateMessage.Contact != nil:
			message.Type = model.MessageContact
			message.Data = fmt.Sprintf("%s - %s %s",
				updateMessage.Contact.PhoneNumber,
				updateMessage.Contact.FirstName,
				updateMessage.Contact.LastName)
			messages <- message
		case updateMessage.Location != nil:
			message.Type = model.MessageLocation
			message.Data = fmt.Sprintf("%v|%v",
				updateMessage.Location.Latitude,
				updateMessage.Location.Longitude)
			messages <- message
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
	case model.MessageContact:
		msg := tgbotapi.NewContact(m.ChatID, m.Contact.PhoneNumber, m.From.FirstName)
		return msg
	case model.MessageLocation:
		msg := tgbotapi.NewLocation(m.ChatID, m.Location.Latitude, m.Location.Longitude)
		return msg
	case model.MessageSticker:
		return tgbotapi.NewStickerUpload(m.ChatID, m.Data)
	case model.MessagePhoto:
		photo := tgbotapi.NewPhotoUpload(m.ChatID, m.Data)
		photo = tgbotapi.PhotoConfig{BaseFile: fileMessage(m, photo.BaseFile)}
		photo.Caption = m.Caption
		return photo
	case model.MessageAudio:
		msg := tgbotapi.NewAudioUpload(m.ChatID, m.Data)
		msg = tgbotapi.AudioConfig{BaseFile: fileMessage(m, msg.BaseFile)}
		return msg
	case model.MessageDocument:
		msg := tgbotapi.NewDocumentUpload(m.ChatID, nil)
		msg = tgbotapi.DocumentConfig{BaseFile: fileMessage(m, msg.BaseFile)}
		return msg
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
	case model.MessageContactButton:
		msg := tgbotapi.NewMessage(m.ChatID, m.Data)
		btn := [][]tgbotapi.KeyboardButton{
			[]tgbotapi.KeyboardButton{
				tgbotapi.NewKeyboardButtonContact(m.ContactButton)}}
		keyboard := tgbotapi.NewReplyKeyboard(btn...)
		keyboard.OneTimeKeyboard = m.OneTimeKeyboard
		msg.ReplyMarkup = keyboard
		if m.Markdown {
			msg.ParseMode = tgbotapi.ModeMarkdown
		}
		return msg
	case model.MessageLocationButton:
		msg := tgbotapi.NewMessage(m.ChatID, m.Data)
		btn := [][]tgbotapi.KeyboardButton{
			[]tgbotapi.KeyboardButton{
				tgbotapi.NewKeyboardButtonLocation(m.LocationButton)}}
		keyboard := tgbotapi.NewReplyKeyboard(btn...)
		keyboard.OneTimeKeyboard = m.OneTimeKeyboard
		msg.ReplyMarkup = keyboard
		if m.Markdown {
			msg.ParseMode = tgbotapi.ModeMarkdown
		}
		return msg
	case model.MessageInlineKeyboard:
		msg := tgbotapi.NewMessage(m.ChatID, m.Data)
		var btns [][]tgbotapi.InlineKeyboardButton
		if m.WithDataInlineButtons {
			btns = inlineDataButtonsFromStrings(m.InlineButtons)
		} else if m.WithURLInlineButtons {
			btns = inlineURLButtonsFromStrings(m.InlineButtons)
		}
		keyboard := tgbotapi.NewInlineKeyboardMarkup(btns...)
		msg.ReplyMarkup = keyboard
		if m.Markdown {
			msg.ParseMode = tgbotapi.ModeMarkdown
		}
		return msg
	}
	return nil
}

func inlineDataButtonsFromStrings(strs []map[string]string) [][]tgbotapi.InlineKeyboardButton {
	btns := make([][]tgbotapi.InlineKeyboardButton, len(strs))
	for i, buttonRow := range strs {
		btnsRow := []tgbotapi.InlineKeyboardButton{}
		for buttonText, buttonData := range buttonRow {
			btnsRow = append(btnsRow, tgbotapi.NewInlineKeyboardButtonData(buttonText, buttonData))
		}
		btns[i] = tgbotapi.NewInlineKeyboardRow(btnsRow...)
	}
	return btns
}

func inlineURLButtonsFromStrings(strs []map[string]string) [][]tgbotapi.InlineKeyboardButton {
	btns := make([][]tgbotapi.InlineKeyboardButton, len(strs))
	for i, buttonRow := range strs {
		btnsRow := []tgbotapi.InlineKeyboardButton{}
		for buttonText, buttonURL := range buttonRow {
			btnsRow = append(btnsRow, tgbotapi.NewInlineKeyboardButtonURL(buttonText, buttonURL))
		}
		btns[i] = tgbotapi.NewInlineKeyboardRow(btnsRow...)
	}
	return btns
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

func fileMessage(m *model.Message, file tgbotapi.BaseFile) tgbotapi.BaseFile {
	if strings.HasPrefix(m.Data, "http") {
		_, err := url.Parse(m.Data)
		if err != nil {
			return file
		}
		file.FileID = m.Data
		file.UseExisting = true
		return file
	}
	return file
}
