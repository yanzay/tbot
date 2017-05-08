package tbot

import "github.com/yanzay/tbot/adapter"

// Server is a telegram bot server. Looks and feels like net/http.
type Server struct {
	bot         *adapter.Bot
	mux         Mux
	middlewares []Middleware
}

type Middleware func(HandlerFunction) HandlerFunction

// NewServer creates new Server with Telegram API Token
// and default /help handler
func NewServer(token string) (*Server, error) {
	tbot, err := adapter.CreateBot(token)
	if err != nil {
		return nil, err
	}

	server := &Server{
		bot: tbot,
		mux: NewDefaultMux(),
	}

	server.HandleFunc("/help", server.HelpHandler)

	return server, nil
}

func (s *Server) AddMiddleware(mid Middleware) {
	s.middlewares = append(s.middlewares, mid)
}

// ListenAndServe starts Server, returns error on failure
func (s *Server) ListenAndServe() error {
	updates, err := s.bot.GetUpdatesChan()
	if err != nil {
		return err
	}
	for update := range updates {
		go s.processMessage(&Message{Message: update})
	}
	return nil
}

// HandleFunc delegates HandleFunc to the current Mux
func (s *Server) HandleFunc(path string, handler HandlerFunction, description ...string) {
	s.mux.HandleFunc(path, handler, description...)
}

// Handle is a shortcut for HandleFunc to reply just with static text,
// "description" is for "/help" handler.
func (s *Server) Handle(path string, reply string, description ...string) {
	f := func(m *Message) {
		m.Reply(reply)
	}
	s.HandleFunc(path, f, description...)
}

func (s *Server) HandleFile(handler HandlerFunction, description ...string) {
	s.mux.HandleFile(handler, description...)
}

// HandleDefault delegates HandleDefault to the current Mux
func (s *Server) HandleDefault(handler HandlerFunction, description ...string) {
	s.mux.HandleDefault(handler, description...)
}
