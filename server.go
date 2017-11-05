package tbot

import (
	"net/http"

	"github.com/yanzay/tbot/internal/adapter"
	"github.com/yanzay/tbot/model"
)

// Server is a telegram bot server. Looks and feels like net/http.
type Server struct {
	bot         adapter.BotAdapter
	mux         Mux
	middlewares []Middleware
	webhookURL  string
	listenAddr  string
}

// Middleware function takes HandlerFunction and returns HandlerFunction.
// Should call it's argument function inside, if needed.
type Middleware func(HandlerFunction) HandlerFunction

var createBot = func(token string, httpClient *http.Client) (adapter.BotAdapter, error) {
	return adapter.CreateBot(token, httpClient)
}

// ServerOption is a functional option for Server
type ServerOption func(*Server)

// WithWebhook returns ServerOption for given Webhook URL and Server address to listen.
// e.g. WithWebook("https://bot.example.com/super/url", "0.0.0.0:8080")
func WithWebhook(url string, addr string) ServerOption {
	return func(s *Server) {
		s.webhookURL = url
		s.listenAddr = addr
	}
}

// WithMux sets custom mux for server. Should satisfy Mux interface.
func WithMux(m Mux) ServerOption {
	return func(s *Server) {
		s.mux = m
	}
}

// NewServer creates new Server with Telegram API Token
// and default /help handler using go default http client
func NewServer(token string, options ...ServerOption) (*Server, error) {
	return NewServerWithClient(token, http.DefaultClient, options...)
}

// NewServerWithClient creates new Server with Telegram API Token
// and default /help handler
func NewServerWithClient(token string, httpClient *http.Client, options ...ServerOption) (*Server, error) {
	tbot, err := createBot(token, httpClient)
	if err != nil {
		return nil, err
	}

	server := &Server{
		bot: tbot,
		mux: NewDefaultMux(),
	}

	for _, option := range options {
		option(server)
	}

	server.HandleFunc("/help", server.HelpHandler)

	return server, nil
}

// AddMiddleware adds new Middleware for server
func (s *Server) AddMiddleware(mid Middleware) {
	s.middlewares = append(s.middlewares, mid)
}

// ListenAndServe starts Server, returns error on failure
func (s *Server) ListenAndServe() error {
	updates, err := s.bot.GetUpdatesChan(s.webhookURL, s.listenAddr)
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

// HandleFile adds file handler for user uploads.
func (s *Server) HandleFile(handler HandlerFunction, description ...string) {
	s.mux.HandleFile(handler, description...)
}

// HandleDefault delegates HandleDefault to the current Mux
func (s *Server) HandleDefault(handler HandlerFunction, description ...string) {
	s.mux.HandleDefault(handler, description...)
}

func (s *Server) SetAlias(route string, aliases ...string) {
	s.mux.SetAlias(route, aliases...)
}

func (s *Server) Send(chatID int64, text string) error {
	return s.bot.Send(&model.Message{Type: model.MessageText, ChatID: chatID, Data: text})
}

func (s *Server) Reset(chatID int64) {
	s.mux.Reset(chatID)
}
