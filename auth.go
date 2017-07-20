package tbot

// NewAuth creates Middleware for white-list based authentication
func NewAuth(whitelist []string) Middleware {
	return func(f HandlerFunction) HandlerFunction {
		return func(m *Message) {
			for _, name := range whitelist {
				if m.From.UserName == name {
					f(m)
					return
				}
			}
			accessDenied(m)
		}
	}
}

func accessDenied(m *Message) {
	m.Replyf("Access denied for user %s", m.From.UserName)
}
