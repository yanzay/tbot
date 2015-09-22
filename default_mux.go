package tbot

import "regexp"

func DefaultMux(handlers map[string]*Handler, path string) (*Handler, MessageVars) {
	for _, handler := range handlers {
		re := regexp.MustCompile(handler.pattern)
		matches := re.FindStringSubmatch(path)

		if len(matches) > 0 {
			messageData := make(map[string]string)
			matches := matches[1:]
			for i, match := range matches {
				messageData[handler.variables[i]] = match
			}
			return handler, messageData
		}
	}
	return nil, nil
}

func parseVariables(pattern string) []string {
	vars := make([]string, 0)
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
