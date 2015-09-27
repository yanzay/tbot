package tbot

import "regexp"

type DefaultMux struct {
	handlers       Handlers
	defaultHandler *Handler
}

func NewDefaultMux() Mux {
	return &DefaultMux{handlers: make(Handlers)}
}

func (dm *DefaultMux) Handlers() Handlers {
	return dm.handlers
}

func (dm *DefaultMux) DefaultHandler() *Handler {
	return dm.defaultHandler
}

func (dm *DefaultMux) Mux(path string) (*Handler, MessageVars) {
	for _, handler := range dm.handlers {
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
	return dm.defaultHandler, nil
}

func (dm *DefaultMux) HandleFunc(path string, handler HandlerFunction, description ...string) {
	dm.handlers[path] = NewHandler(handler, path, description...)
}

func (dm *DefaultMux) Handle(path string, reply string, description ...string) {
	f := func(m Message) {
		m.Reply(reply)
	}
	dm.HandleFunc(path, f, description...)
}

func (dm *DefaultMux) HandleDefault(handler HandlerFunction, description ...string) {
	dm.defaultHandler = NewHandler(handler, "", description...)
}
