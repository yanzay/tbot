package tbot

import (
	"fmt"
	"strings"
)

func (s *Server) HelpHandler(m Message) {
	handlerNames := make([]string, 0)
	for handlerName, handler := range s.handlers {
		var line string
		line = handlerName
		if handler.description != "" {
			line = fmt.Sprintf("%s - %s", line, handler.description)
		}
		handlerNames = append(handlerNames, line)
	}
	if s.defaultHandler != nil && s.defaultHandler.description != "" {
		defaultLine := fmt.Sprintf("* - %s", s.defaultHandler.description)
		handlerNames = append(handlerNames, defaultLine)
	}
	reply := strings.Join(handlerNames, "\n")
	m.Reply(reply)
}
