package tbot

func NewAuth(whitelist []string) Middleware {
	return func(f HandlerFunction) HandlerFunction {
		return func(m *Message) {
			for _, name := range whitelist {
				if m.From == name {
					f(m)
					return
				}
			}
			AccessDenied(m)
		}
	}
}

func AccessDenied(m *Message) {
	m.Replyf("Access denied for user %s", m.From)
}
