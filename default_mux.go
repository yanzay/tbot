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
