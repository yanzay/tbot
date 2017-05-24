package tbot

import "strings"

const (
	RouteBack    = "/<..>"
	RouteRoot    = "/<root>"
	RouteRefresh = "/<.>"
)

// RouterMux is a tree-route multiplexer
type RouterMux struct {
	handlers       Handlers
	fileHandler    *Handler
	defaultHandler *Handler
	storage        SessionStorage
	aliases        map[string]string
}

// NewRouterMux creates new RouterMux
// Takes SessionStorage to store users' sessions state
func NewRouterMux(storage SessionStorage) Mux {
	return &RouterMux{
		handlers: make(Handlers),
		storage:  storage,
		aliases:  make(map[string]string),
	}
}

// Handlers returns list of handlers currently presented in mux
func (rm *RouterMux) Handlers() Handlers {
	return rm.handlers
}

// DefaultHandler returns default handler, nil if it's not set
func (rm *RouterMux) DefaultHandler() *Handler {
	return rm.defaultHandler
}

func (rm *RouterMux) FileHandler() *Handler {
	return rm.fileHandler
}

// Mux takes message content and returns corresponding handler
func (rm *RouterMux) Mux(msg *Message) (*Handler, MessageVars) {
	state := rm.storage.Get(msg.ChatID)
	if state == "" {
		state = RouteRoot
	}
	route := msg.Data
	if _, ok := rm.aliases[route]; ok {
		route = rm.aliases[route]
	}
	switch route {
	case RouteBack:
		state = back(state)
	case RouteRoot:
		state = RouteRoot
	case RouteRefresh:
	default:
		if rm.handlers[state+route] != nil {
			state += route
		} else if rm.handlers[back(state)+route] != nil {
			state = back(state) + route
		}
	}
	if rm.handlers[state] == nil {
		return rm.defaultHandler, nil
	}
	rm.storage.Set(msg.ChatID, state)
	return rm.handlers[state], MessageVars{}
}

func back(route string) string {
	routes := strings.Split(route, "/")
	return strings.Join(routes[:len(routes)-1], "/")
}

// HandleFunc adds new handler function to mux.
func (rm *RouterMux) HandleFunc(path string, handler HandlerFunction, description ...string) {
	if path != RouteRoot {
		path = RouteRoot + path
	}
	rm.handlers[path] = NewHandler(handler, path, description...)
}

// SetAlias sets aliases for specified route.
func (rm *RouterMux) SetAlias(route string, aliases ...string) {
	for _, alias := range aliases {
		rm.aliases[alias] = route
	}
}

// HandleDefault adds new default handler, when nothing matches with message,
func (rm *RouterMux) HandleDefault(handler HandlerFunction, description ...string) {
	rm.defaultHandler = NewHandler(handler, "", description...)
}

func (rm *RouterMux) HandleFile(handler HandlerFunction, description ...string) {
}

func (rm *RouterMux) Reset(chatID int64) {
	rm.storage.Reset(chatID)
}
