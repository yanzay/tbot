package tbot

type HandlerFunction func(Message)
type Handlers map[string]*Handler

type Mux interface {
	Mux(string) (*Handler, MessageVars)
	HandleFunc(string, HandlerFunction, ...string)
	Handle(string, string, ...string)
	HandleDefault(HandlerFunction, ...string)

	Handlers() Handlers
	DefaultHandler() *Handler
}
