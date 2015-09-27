package tbot

import (
	"log"
	"time"

	"github.com/tucnak/telebot"
)

type Server struct {
	bot *telebot.Bot
	mux Mux
}

func NewServer(token string) (*Server, error) {
	tbot, err := telebot.NewBot(token)
	if err != nil {
		return nil, err
	}

	server := &Server{
		bot: tbot,
		mux: NewDefaultMux(),
	}

	return server, nil
}

func (s *Server) ListenAndServe() {
	messages := s.listenMessages(3 * time.Second)
	for message := range messages {
		go s.processMessage(message)
	}
}

func (s *Server) processMessage(message telebot.Message) {
	log.Printf("[TBot] %s %s: %s", message.Sender.FirstName, message.Sender.LastName, message.Text)
	handler, data := s.mux.Mux(message.Text)
	if handler == nil {
		return
	}
	m := Message{message, data, make(chan *ReplyMessage), make(chan struct{})}
	go func() {
		handler.f(m)
		m.close <- struct{}{}
	}()
	for {
		select {
		case reply := <-m.replies:
			switch reply.messageType {
			case MessageText:
				s.bot.SendMessage(message.Chat, reply.Text, nil)
			case MessageSticker:
				s.bot.SendSticker(message.Chat, &reply.Sticker, nil)
			case MessagePhoto:
				s.bot.SendPhoto(message.Chat, reply.photo, nil)
			case MessageAudio:
				s.bot.SendAudio(message.Chat, reply.audio, nil)
			case MessageDocument:
				s.bot.SendDocument(message.Chat, reply.document, nil)
			}
		case <-m.close:
			return
		}
	}
}

func (s *Server) listenMessages(interval time.Duration) <-chan telebot.Message {
	messages := make(chan telebot.Message)
	s.bot.Listen(messages, interval)
	return messages
}
