package tbot

type HandlerFunction func(Message)
type Mux func(map[string]*Handler, string) (*Handler, MessageVars)
