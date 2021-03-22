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

const apiBaseURL = "https://api.telegram.org"

// Server will connect and serve all updates from Telegram
type Server struct {
	webhookURL      string
	webhookHandler  func(w http.ResponseWriter, r *http.Request)
	useCustomServer bool
	listenAddr      string
	baseURL         string
	httpClient      *http.Client
	client          *Client
	token           string
	logger          Logger
	stop            chan struct{}
	updates         chan *Update
	updatesParams   url.Values
	bufferSize      int
	nextOffset      int

	messageHandlers        []messageHandler
	editMessageHandler     handlerFunc
	channelPostHandler     handlerFunc
	editChannelPostHandler handlerFunc
	inlineQueryHandler     func(*InlineQuery)
	inlineResultHandler    func(*ChosenInlineResult)
	callbackHandler        func(*CallbackQuery)
	shippingHandler        func(*ShippingQuery)
	preCheckoutHandler     func(*PreCheckoutQuery)
	pollHandler            func(*Poll)
	pollAnswerHandler      func(*PollAnswer)

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
	WithWebhook(url, addr string)
	WithWebhookForCustomServer(url string)
	WithHTTPClient(client *http.Client)
	WithBaseURL(baseURL string)
*/
func New(token string, options ...ServerOption) *Server {
	s := &Server{
		httpClient:      http.DefaultClient,
		token:           token,
		logger:          nopLogger{},
		baseURL:         apiBaseURL,
		useCustomServer: false,

		editMessageHandler:     func(*Message) {},
		channelPostHandler:     func(*Message) {},
		editChannelPostHandler: func(*Message) {},
		inlineQueryHandler:     func(*InlineQuery) {},
		inlineResultHandler:    func(*ChosenInlineResult) {},
		callbackHandler:        func(*CallbackQuery) {},
		shippingHandler:        func(*ShippingQuery) {},
		preCheckoutHandler:     func(*PreCheckoutQuery) {},
		pollHandler:            func(*Poll) {},
		pollAnswerHandler:      func(*PollAnswer) {},

		stop: make(chan struct{}, 0),
	}
	for _, opt := range options {
		opt(s)
	}

	s.updates = make(chan *Update, s.bufferSize)
	s.webhookHandler = func(w http.ResponseWriter, r *http.Request) {
		up := &Update{}
		err := json.NewDecoder(r.Body).Decode(up)
		if err != nil {
			s.logger.Errorf("unable to decode update: %v", err)
			return
		}
		s.updates <- up
	}

	// bot, err :=  tgbotapi.NewBotAPIWithClient(token, s.httpClient)
	s.client = NewClient(token, s.httpClient, s.baseURL)
	return s
}

// WithWebhook returns ServerOption for given Webhook URL and Server address to listen.
// e.g. WithWebhook("https://bot.example.com/super/url", "0.0.0.0:8080")
func WithWebhook(url, addr string) ServerOption {
	return func(s *Server) {
		s.webhookURL = url
		s.listenAddr = addr
	}
}

// WithWebhookForCustomServer sets the useCustomServer to true in order to handle the s.webhookHandler
// with an custom server implementation
func WithWebhookForCustomServer(url string) ServerOption {
	return func(s *Server) {
		s.bufferSize = 5
		s.webhookURL = url
		s.useCustomServer = true
	}
}

// WithBaseURL sets custom apiBaseURL for server.
// It may be necessary to run the server in some countries
func WithBaseURL(baseURL string) ServerOption {
	return func(s *Server) {
		s.baseURL = baseURL
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
	err := s.initializeUpdates()
	if err != nil {
		return err
	}
	for {
		select {
		case update := <-s.updates:
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
				case update.Poll != nil:
					s.pollHandler(update.Poll)
				case update.PollAnswer != nil:
					s.pollAnswerHandler(update.PollAnswer)
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

// GetWebhookHandler returns the webhook handler
func (s *Server) GetWebhookHandler() func(w http.ResponseWriter, r *http.Request) {
	return s.webhookHandler
}

func (s *Server) initializeUpdates() error {
	if s.useCustomServer || (s.webhookURL != "" && s.listenAddr != "") {
		return s.listenUpdates()
	}
	s.client.deleteWebhook()
	return s.longPoolUpdates()
}

func (s *Server) listenUpdates() error {
	err := s.client.setWebhook(s.webhookURL)
	if err != nil {
		return fmt.Errorf("unable to set webhook: %v", err)
	}

	if !s.useCustomServer {
		// Start standalone server for webhookHandler
		l, err := net.Listen("tcp", s.listenAddr)
		if err != nil {
			return err
		}
		go http.Serve(l, http.HandlerFunc(s.webhookHandler))
	}

	return nil
}

func (s *Server) longPoolUpdates() error {
	s.logger.Debugf("fetching updates...")
	endpoint := fmt.Sprintf("%s/bot%s/%s", s.baseURL, s.token, "getUpdates")
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}
	params := s.updatesParams
	if params == nil {
		params = url.Values{}
	}
	params.Set("timeout", fmt.Sprint(3600))
	req.URL.RawQuery = params.Encode()
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
				s.updates <- up
			}
		}
	}()
	return nil
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

// HandlePollUpdate set handler for anonymous poll updates
func (s *Server) HandlePollUpdate(handler func(*Poll)) {
	s.pollHandler = handler
}

// HandlePollAnswer set handler for non-anonymous poll updates
func (s *Server) HandlePollAnswer(handler func(*PollAnswer)) {
	s.pollAnswerHandler = handler
}

func (s *Server) handleMessage(msg *Message) {
	for _, handler := range s.messageHandlers {
		if handler.rx.MatchString(msg.Text) {
			handler.f(msg)
			return
		}
	}
}
