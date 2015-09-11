package tbot

type Handler func(Message)
type Mux func(map[string]Handler, string) (Handler, MessageVars)
