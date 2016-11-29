package tbot

import (
	"fmt"
	"strings"
)

// HelpHandler is a default handler for /help,
// shows available commands and their description
func (s *Server) HelpHandler(m Message) {
	var handlerNames []string
	for handlerName, handler := range s.mux.Handlers() {
		var line string
		line = handlerName
		if handler.description != "" {
			line = fmt.Sprintf("%s - %s", line, handler.description)
		}
		handlerNames = append(handlerNames, line)
	}
	if s.mux.DefaultHandler() != nil && s.mux.DefaultHandler().description != "" {
		defaultLine := fmt.Sprintf("* - %s", s.mux.DefaultHandler().description)
		handlerNames = append(handlerNames, defaultLine)
	}
	reply := strings.Join(handlerNames, "\n")
	m.Reply(reply)
}
