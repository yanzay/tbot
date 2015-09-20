package tbot

type HandlerFunction func(Message)

type Handler struct {
	f           HandlerFunction
	description string
}

func NewHandler(f func(Message), description ...string) *Handler {
	handler := &Handler{}
	handler.f = f
	if len(description) > 0 {
		handler.description = description[0]
	}
	return handler
}

type Mux func(map[string]*Handler, string) (*Handler, MessageVars)
