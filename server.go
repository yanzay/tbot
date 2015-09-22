package tbot

import (
	"log"
	"time"

	"github.com/tucnak/telebot"
)

type Server struct {
	bot            *telebot.Bot
	mux            Mux
	handlers       map[string]*Handler
	defaultHandler *Handler
}

func NewServer(token string) (*Server, error) {
	tbot, err := telebot.NewBot(token)
	if err != nil {
		return nil, err
	}

	server := &Server{
		bot:      tbot,
		handlers: make(map[string]*Handler),
		mux:      DefaultMux,
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
	handler, data := s.mux(s.handlers, message.Text)
	if handler == nil {
		handler = s.defaultHandler
	}
	if handler != nil {
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
				}
			case <-m.close:
				return
			}
		}
	}
}

func (s *Server) HandleFunc(path string, handler HandlerFunction, description ...string) {
	s.handlers[path] = NewHandler(handler, path, description...)
}

func (s *Server) Handle(path string, reply string, description ...string) {
	f := func(m Message) {
		m.Reply(reply)
	}
	s.HandleFunc(path, f, description...)
}

func (b *Server) HandleDefault(handler HandlerFunction, description ...string) {
	b.defaultHandler = NewHandler(handler, "", description...)
}

func (s *Server) listenMessages(interval time.Duration) <-chan telebot.Message {
	messages := make(chan telebot.Message)
	s.bot.Listen(messages, interval)
	return messages
}
