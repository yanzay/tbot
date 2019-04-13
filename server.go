package tbot

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

var (
	apiBaseURL = "https://api.telegram.org"
)

// Server will connect and serve all updates from Telegram
type Server struct {
	webhookURL    string
	listenAddr    string
	httpClient    *http.Client
	client        *Client
	token         string
	logger        Logger
	stop          chan struct{}
	updatesParams url.Values
	bufferSize    int
	nextOffset    int

	messageHandlers        []messageHandler
	editMessageHandler     handlerFunc
	channelPostHandler     handlerFunc
	editChannelPostHandler handlerFunc
	inlineQueryHandler     func(*InlineQuery)
	inlineResultHandler    func(*ChosenInlineResult)
	callbackHandler        func(*CallbackQuery)
	shippingHandler        func(*ShippingQuery)
	preCheckoutHandler     func(*PreCheckoutQuery)

	middlewares []Middleware
}

// UpdateHandler is a function for middlewares
type UpdateHandler func(*Update)

// Middleware is a middleware for updates
type Middleware func(UpdateHandler) UpdateHandler

// ServerOption type for additional Server options
type ServerOption func(*Server)

type handlerFunc func(*Message)

type messageHandler struct {
	rx *regexp.Regexp
	f  handlerFunc
}

/*
New creates new Server. Available options:
	WithWebook(url, addr string)
	WithHTTPClient(client *http.Client)
*/
func New(token string, options ...ServerOption) *Server {
	s := &Server{
		httpClient: http.DefaultClient,
		token:      token,
		logger:     nopLogger{},

		editMessageHandler:     func(*Message) {},
		channelPostHandler:     func(*Message) {},
		editChannelPostHandler: func(*Message) {},
		inlineQueryHandler:     func(*InlineQuery) {},
		inlineResultHandler:    func(*ChosenInlineResult) {},
		callbackHandler:        func(*CallbackQuery) {},
		shippingHandler:        func(*ShippingQuery) {},
		preCheckoutHandler:     func(*PreCheckoutQuery) {},
	}
	for _, opt := range options {
		opt(s)
	}
	// bot, err :=  tgbotapi.NewBotAPIWithClient(token, s.httpClient)
	s.client = NewClient(token, s.httpClient, apiBaseURL)
	return s
}

// WithWebhook returns ServerOption for given Webhook URL and Server address to listen.
// e.g. WithWebook("https://bot.example.com/super/url", "0.0.0.0:8080")
func WithWebhook(url, addr string) ServerOption {
	return func(s *Server) {
		s.webhookURL = url
		s.listenAddr = addr
	}
}

// WithHTTPClient sets custom http client for server.
func WithHTTPClient(client *http.Client) ServerOption {
	return func(s *Server) {
		s.httpClient = client
	}
}

// WithLogger sets logger for tbot
func WithLogger(logger Logger) ServerOption {
	return func(s *Server) {
		s.logger = logger
	}
}

// Use adds middleware to server
func (s *Server) Use(m Middleware) {
	s.middlewares = append(s.middlewares, m)
}

// Start listening for updates
func (s *Server) Start() error {
	if len(s.token) == 0 {
		return fmt.Errorf("token is empty")
	}
	updates, err := s.getUpdates()
	if err != nil {
		return err
	}
	for {
		select {
		case update := <-updates:
			handleUpdate := func(update *Update) {
				switch {
				case update.Message != nil:
					s.handleMessage(update.Message)
				case update.EditedMessage != nil:
					s.editMessageHandler(update.EditedMessage)
				case update.ChannelPost != nil:
					s.channelPostHandler(update.ChannelPost)
				case update.EditedChannelPost != nil:
					s.editChannelPostHandler(update.EditedChannelPost)
				case update.InlineQuery != nil:
					s.inlineQueryHandler(update.InlineQuery)
				case update.ChosenInlineResult != nil:
					s.inlineResultHandler(update.ChosenInlineResult)
				case update.CallbackQuery != nil:
					s.callbackHandler(update.CallbackQuery)
				case update.ShippingQuery != nil:
					s.shippingHandler(update.ShippingQuery)
				case update.PreCheckoutQuery != nil:
					s.preCheckoutHandler(update.PreCheckoutQuery)
				}
			}
			var f = handleUpdate
			for i := len(s.middlewares) - 1; i >= 0; i-- {
				f = s.middlewares[i](f)
			}
			go f(update)
		case <-s.stop:
			return nil
		}
	}
}

// Client returns Telegram API Client
func (s *Server) Client() *Client {
	return s.client
}

// Stop listening for updates
func (s *Server) Stop() {
	s.stop <- struct{}{}
}

func (s *Server) getUpdates() (chan *Update, error) {
	if s.webhookURL != "" && s.listenAddr != "" {
		return s.listenUpdates()
	}
	s.client.deleteWebhook()
	return s.longPoolUpdates()
}

func (s *Server) listenUpdates() (chan *Update, error) {
	err := s.client.setWebhook(s.webhookURL)
	if err != nil {
		return nil, fmt.Errorf("unable to set webhook: %v", err)
	}
	updates := make(chan *Update)
	handler := func(w http.ResponseWriter, r *http.Request) {
		up := &Update{}
		err := json.NewDecoder(r.Body).Decode(up)
		if err != nil {
			s.logger.Errorf("unable to decode update: %v", err)
			return
		}
		updates <- up
	}
	l, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return nil, err
	}
	go http.Serve(l, http.HandlerFunc(handler))
	return updates, nil
}

func (s *Server) longPoolUpdates() (chan *Update, error) {
	s.logger.Debugf("fetching updates...")
	endpoint := fmt.Sprintf("%s/bot%s/%s", apiBaseURL, s.token, "getUpdates")
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	params := s.updatesParams
	if params == nil {
		params = url.Values{}
	}
	params.Set("timeout", fmt.Sprint(3600))
	req.URL.RawQuery = params.Encode()
	updates := make(chan *Update, s.bufferSize)
	go func() {
		for {
			params.Set("offset", fmt.Sprint(s.nextOffset))
			req.URL.RawQuery = params.Encode()
			resp, err := s.httpClient.Do(req)
			if err != nil {
				s.logger.Errorf("unable to perform request: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}
			var updatesResp *struct {
				OK          bool      `json:"ok"`
				Result      []*Update `json:"result"`
				Description string    `json:"description"`
			}
			err = json.NewDecoder(resp.Body).Decode(&updatesResp)
			if err != nil {
				s.logger.Errorf("unable to decode response: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}
			err = resp.Body.Close()
			if err != nil {
				s.logger.Errorf("unable to close response body: %v", err)
			}
			if !updatesResp.OK {
				s.logger.Errorf("updates query fail: %s", updatesResp.Description)
				time.Sleep(1 * time.Second)
				continue
			}
			for _, up := range updatesResp.Result {
				s.nextOffset = up.UpdateID + 1
				updates <- up
			}
		}
	}()
	return updates, nil
}

// HandleMessage sets handler for incoming messages
func (s *Server) HandleMessage(pattern string, handler func(*Message)) {
	rx := regexp.MustCompile(pattern)
	s.messageHandlers = append(s.messageHandlers, messageHandler{rx: rx, f: handler})
}

// HandleEditedMessage set handler for incoming edited messages
func (s *Server) HandleEditedMessage(handler func(*Message)) {
	s.editMessageHandler = handler
}

// HandleChannelPost set handler for incoming channel post
func (s *Server) HandleChannelPost(handler func(*Message)) {
	s.channelPostHandler = handler
}

// HandleEditChannelPost set handler for incoming edited channel post
func (s *Server) HandleEditChannelPost(handler func(*Message)) {
	s.editChannelPostHandler = handler
}

// HandleInlineQuery set handler for inline queries
func (s *Server) HandleInlineQuery(handler func(*InlineQuery)) {
	s.inlineQueryHandler = handler
}

// HandleInlineResult set inline result handler
func (s *Server) HandleInlineResult(handler func(*ChosenInlineResult)) {
	s.inlineResultHandler = handler
}

// HandleCallback set handler for inline buttons
func (s *Server) HandleCallback(handler func(*CallbackQuery)) {
	s.callbackHandler = handler
}

// HandleShipping set handler for shipping queries
func (s *Server) HandleShipping(handler func(*ShippingQuery)) {
	s.shippingHandler = handler
}

// HandlePreCheckout set handler for pre-checkout queries
func (s *Server) HandlePreCheckout(handler func(*PreCheckoutQuery)) {
	s.preCheckoutHandler = handler
}

func (s *Server) handleMessage(msg *Message) {
	for _, handler := range s.messageHandlers {
		if handler.rx.MatchString(msg.Text) {
			handler.f(msg)
			return
		}
	}
}
