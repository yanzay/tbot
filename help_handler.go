package tbot

import "strings"

func (s *Server) HelpHandler(m Message) {
	handlerNames := make([]string, 0)
	for handlerName, _ := range s.handlers {
		handlerNames = append(handlerNames, handlerName)
	}
	reply := strings.Join(handlerNames, "\n")
	m.Reply(reply)
}
