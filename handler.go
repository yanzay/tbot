package tbot

import "fmt"

type Handler struct {
	f           HandlerFunction
	description string
	pattern     string
	variables   []string
}

func NewHandler(f func(Message), path string, description ...string) *Handler {
	handler := &Handler{f: f}
	handler.variables = parseVariables(path)
	handler.pattern = fmt.Sprintf("^%s$", replaceVariables(path))
	if len(description) > 0 {
		handler.description = description[0]
	}
	return handler
}
