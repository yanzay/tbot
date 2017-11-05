package tbot

import (
	"strings"

	"github.com/variar/tbot/model"
	"github.com/yanzay/log"
)

func (s *Server) processMessage(message *Message) {
	if message == nil {
		log.Debugf("nil message, return")
		return
	}
	handler, data := s.chooseHandler(message)
	if handler == nil {
		log.Debugf("nil handler, return")
		return
	}
	f := handler.f
	for _, mid := range s.middlewares {
		f = mid(f)
	}
	go s.messageLoop(message.Replies)
	message.Vars = data
	f(message)
	close(message.Replies)
}

func (s *Server) chooseHandler(message *Message) (*Handler, MessageVars) {
	var handler *Handler
	var data MessageVars
	if message.Type == model.MessageDocument {
		handler = s.mux.FileHandler()
		data = map[string]string{"url": message.Data}
	} else {
		message.Data = s.trimBotName(message.Data)
		handler, data = s.mux.Mux(message)
	}

	return handler, data
}

func (s *Server) messageLoop(replies <-chan *model.Message) {
	for reply := range replies {
		err := s.bot.Send(reply)
		if err != nil {
			log.Printf("Error sending message: %q", err)
		}
	}
}

func (s *Server) trimBotName(message string) string {
	parts := strings.SplitN(message, " ", 2)
	command := parts[0]
	command = strings.TrimSuffix(command, "@"+s.bot.GetUserName())
	command = strings.TrimSuffix(command, "@"+s.bot.GetFirstName())
	parts[0] = command
	return strings.Join(parts, " ")
}
