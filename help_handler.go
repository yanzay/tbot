package tbot

import (
	"fmt"
	"strings"
)

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
