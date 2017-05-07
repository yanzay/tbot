package tbot

import (
	"log"
	"strings"

	"github.com/yanzay/tbot/adapter"
)

func (s *Server) processMessage(message *adapter.Message) {
	if message == nil {
		return
	}
	handler, data := s.chooseHandler(message)
	if handler == nil {
		return
	}
	f := handler.f
	for _, mid := range s.middlewares {
		f = mid(f)
	}
	go s.messageLoop(message.Replies)
	f(&Message{Message: message})
	close(message.Replies)
}

func (s *Server) chooseHandler(message *adapter.Message) (*Handler, MessageVars) {
	var handler *Handler
	var data MessageVars
	if message.Type == adapter.MessageDocument {
		handler = s.mux.FileHandler()
		data = map[string]string{"url": message.Data}
	} else {
		message.Data = s.trimBotName(message.Data)
		handler, data = s.mux.Mux(message.Data)
	}

	return handler, data
}

func (s *Server) messageLoop(replies <-chan *adapter.Message) {
	for reply := range replies {
		err := s.dispatchMessage(reply)
		if err != nil {
			log.Printf("Error dispatching message: %q", err)
		}
	}
}

func (s *Server) dispatchMessage(reply *adapter.Message) error {
	return s.bot.Send(reply)
}

func (s *Server) trimBotName(message string) string {
	parts := strings.SplitN(message, " ", 2)
	command := parts[0]
	command = strings.TrimSuffix(command, "@"+s.bot.GetUserName())
	command = strings.TrimSuffix(command, "@"+s.bot.GetFirstName())
	parts[0] = command
	return strings.Join(parts, " ")
}
