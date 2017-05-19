package tbot

import (
	"fmt"
	"regexp"
)

// Handler is a struct that represents any message handler
// with handler function, description, pattern and parsed variables
type Handler struct {
	f           HandlerFunction
	description string
	pattern     string
	variables   []string
}

// NewHandler creates new handler and returns it
func NewHandler(f func(*Message), path string, description ...string) *Handler {
	handler := &Handler{f: f}
	handler.variables, handler.pattern = parse(path)
	if len(description) > 0 {
		handler.description = description[0]
	}
	return handler
}

func parse(template string) ([]string, string) {
	vars := parseVariables(template)
	pattern := fmt.Sprintf("^%s$", replaceVariables(template))
	return vars, pattern
}

func parseVariables(pattern string) []string {
	var vars []string
	re := regexp.MustCompile("{([A-Za-z0-9_]*)}")
	matches := re.FindAllStringSubmatch(pattern, -1)
	for _, match := range matches {
		if len(match) > 0 {
			vars = append(vars, match[1])
		}
	}
	return vars
}

func replaceVariables(pattern string) string {
	re := regexp.MustCompile("{[A-Za-z0-9_]*}")
	return re.ReplaceAllString(pattern, "(.*)")
}
