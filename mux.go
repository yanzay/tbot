package tbot

// HandlerFunction is a function that can process incoming messages
type HandlerFunction func(Message)

// Handlers is a lookup table of handlers,
// key - string pattern
// value - Handler
type Handlers map[string]*Handler

// Mux interface represents message multiplexer
type Mux interface {
	Mux(string) (*Handler, MessageVars)
	HandleFunc(string, HandlerFunction, ...string)
	HandleDefault(HandlerFunction, ...string)

	Handlers() Handlers
	DefaultHandler() *Handler
}
