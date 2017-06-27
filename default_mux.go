package tbot

import "regexp"

// DefaultMux is a default multiplexer,
// supports parametrized commands.
// Parameters should be enclosed with curly brackets,
// like in "/say {hi}" - "hi" is a parameter.
type DefaultMux struct {
	handlers       Handlers
	fileHandler    *Handler
	defaultHandler *Handler
	aliases        map[string]string
}

// NewDefaultMux creates new DefaultMux
func NewDefaultMux() Mux {
	return &DefaultMux{
		handlers: make(Handlers),
		aliases:  make(map[string]string),
	}
}

// Handlers returns list of handlers currently presented in mux
func (dm *DefaultMux) Handlers() Handlers {
	return dm.handlers
}

// DefaultHandler returns default handler, nil if it's not set
func (dm *DefaultMux) DefaultHandler() *Handler {
	return dm.defaultHandler
}

// FileHandler returns file handler, nil if there is no file handler
func (dm *DefaultMux) FileHandler() *Handler {
	return dm.fileHandler
}

// Mux takes message content and returns corresponding handler
// and parsed vars from message
func (dm *DefaultMux) Mux(msg *Message) (*Handler, MessageVars) {
	path := msg.Data
	for _, handler := range dm.handlers {
		re := regexp.MustCompile(handler.pattern)
		matches := re.FindStringSubmatch(path)

		if len(matches) > 0 {
			messageData := make(map[string]string)
			matches = matches[1:]
			for i, match := range matches {
				messageData[handler.variables[i]] = match
			}
			return handler, messageData
		}
	}
	return dm.defaultHandler, nil
}

func (dm *DefaultMux) SetAlias(route string, aliases ...string) {
	for _, alias := range aliases {
		dm.aliases[alias] = route
	}
}

// HandleFunc adds new handler function to mux, "description" is for "/help" handler.
func (dm *DefaultMux) HandleFunc(path string, handler HandlerFunction, description ...string) {
	dm.handlers[path] = NewHandler(handler, path, description...)
}

// HandleDefault adds new default handler, when nothing matches with message,
// "description" is for "/help" handler.
func (dm *DefaultMux) HandleDefault(handler HandlerFunction, description ...string) {
	dm.defaultHandler = NewHandler(handler, "", description...)
}

// HandleFile adds file handler. When the user uploads Document, Download method will be available for *Message.
func (dm *DefaultMux) HandleFile(handler HandlerFunction, description ...string) {
	dm.fileHandler = NewHandler(handler, "", description...)
}

func (db *DefaultMux) Reset(int64) {
}
