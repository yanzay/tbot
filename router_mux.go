package tbot

import "strings"

const (
	RouteBack    = "/<..>"
	RouteRoot    = "/</>"
	RouteRefresh = "/<.>"
)

// DefaultMux is a default multiplexer,
// supports parametrized commands.
// Parameters should be enclosed with curly brackets,
// like in "/say {hi}" - "hi" is a parameter.
type RouterMux struct {
	handlers       Handlers
	fileHandler    *Handler
	defaultHandler *Handler
	storage        SessionStorage
}

// NewDefaultMux creates new DefaultMux
func NewRouterMux(storage SessionStorage) Mux {
	return &RouterMux{
		handlers: make(Handlers),
		storage:  storage,
	}
}

// Handlers returns list of handlers currently presented in mux
func (rm *RouterMux) Handlers() Handlers {
	return Handlers{}
}

// DefaultHandler returns default handler, nil if it's not set
func (rm *RouterMux) DefaultHandler() *Handler {
	return rm.defaultHandler
}

func (rm *RouterMux) FileHandler() *Handler {
	return rm.fileHandler
}

// Mux takes message content and returns corresponding handler
// and parsed vars from message
func (rm *RouterMux) Mux(msg *Message) (*Handler, MessageVars) {
	state := rm.storage.Get(msg.ChatID)
	switch msg.Data {
	case RouteBack:
		state = back(state)
	case RouteRoot:
		state = root(state)
	case RouteRefresh:
	default:
		state += msg.Data
	}
	rm.storage.Set(msg.ChatID, state)
	return rm.handlers[state], MessageVars{}
}

func back(route string) string {
	routes := strings.Split(route, "/")
	return strings.Join(routes[:len(routes)-1], "/")
}

func root(route string) string {
	routes := strings.SplitN(route, "/", 3)
	return strings.Join(routes[:2], "/")
}

// HandleFunc adds new handler function to mux, "description" is for "/help" handler.
func (rm *RouterMux) HandleFunc(path string, handler HandlerFunction, description ...string) {
	rm.handlers[path] = NewHandler(handler, path, description...)
}

// HandleDefault adds new default handler, when nothing matches with message,
// "description" is for "/help" handler.
func (rm *RouterMux) HandleDefault(handler HandlerFunction, description ...string) {
}

func (rm *RouterMux) HandleFile(handler HandlerFunction, description ...string) {
}
