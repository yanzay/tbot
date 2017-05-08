package adapter

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Bot struct {
	tbot *tgbotapi.BotAPI
}

func CreateBot(token string) (*Bot, error) {
	tbot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &Bot{tbot: tbot}, nil
}

func (b *Bot) GetFileDirectURL(fileID string) (string, error) {
	return b.tbot.GetFileDirectURL(fileID)
}

//func (b *Bot) Send(c tgbotapi.Chattable) error {
func (b *Bot) Send(m *Message) error {
	c := chattableFromMessage(m)
	_, err := b.tbot.Send(c)
	return err
}

func (b *Bot) GetUserName() string {
	return b.tbot.Self.UserName
}

func (b *Bot) GetFirstName() string {
	return b.tbot.Self.FirstName
}

func (b *Bot) GetUpdatesChan() (<-chan *Message, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := b.tbot.GetUpdatesChan(u)
	if err != nil {
		return nil, err
	}
	messages := make(chan *Message)
	go b.adaptUpdates(updates, messages)
	return messages, nil
}

func (b *Bot) adaptUpdates(updates <-chan tgbotapi.Update, messages chan<- *Message) {
	var err error
	for update := range updates {
		message := &Message{Replies: make(chan *Message), From: update.Message.From.UserName}
		switch {
		case update.Message.Document != nil:
			message.Data, err = b.GetFileDirectURL(update.Message.Document.FileID)
			if err != nil {
				log.Println(err)
				continue
			}
			message.Type = MessageDocument
			messages <- message
		case update.Message.Text != "":
			messages <- &Message{Type: MessageText, Data: update.Message.Text}
		}
	}
}

func chattableFromMessage(m *Message) tgbotapi.Chattable {
	switch m.Type {
	case MessageText:
		return tgbotapi.NewMessage(m.ChatID, m.Data)
	case MessageSticker:
		return tgbotapi.NewStickerUpload(m.ChatID, m.Data)
	}
	return nil
}
