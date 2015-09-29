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

func (s *Server) HandleFunc(path string, handler HandlerFunction, description ...string) {
	s.mux.HandleFunc(path, handler, description...)
}

func (s *Server) Handle(path string, reply string, description ...string) {
	s.mux.Handle(path, reply, description...)
}

func (s *Server) HandleDefault(handler HandlerFunction, description ...string) {
	s.mux.HandleDefault(handler, description...)
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
			err := s.dispatchMessage(message.Chat, reply)
			if err != nil {
				log.Println(err)
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

func (s *Server) dispatchMessage(chat telebot.User, reply *ReplyMessage) error {
	var err error
	switch reply.messageType {
	case MessageText:
		err = s.bot.SendMessage(chat, reply.Text, nil)
	case MessageSticker:
		err = s.bot.SendSticker(chat, &reply.Sticker, nil)
	case MessagePhoto:
		err = s.bot.SendPhoto(chat, reply.photo, nil)
	case MessageAudio:
		err = s.bot.SendAudio(chat, reply.audio, nil)
	case MessageDocument:
		err = s.bot.SendDocument(chat, reply.document, nil)
	}
	return err
}
