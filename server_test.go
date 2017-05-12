package tbot

import "testing"

const (
	TestToken    = "TEST:TOKEN"
	InvalidToken = "invalid"
)

func TestNewServerSuccess(t *testing.T) {
	server, err := NewServer(TestToken)
	if err != nil {
		t.Errorf("Error creating server: %s", err)
	}
	if server == nil {
		t.Error("Server is nil")
	}
	if server.mux == nil {
		t.Error("Server mux is nil")
	}
}

func TestNewServerFail(t *testing.T) {
	server, err := NewServer(InvalidToken)
	if err == nil {
		t.Error("Invalid token should return error")
	}
	if server != nil {
		t.Error("Invalid token should return nil server")
	}
}

func TestNewServerWithWebhook(t *testing.T) {
	url := "https://some.url"
	addr := "0.0.0.0:8013"
	server, err := NewServer(TestToken, WithWebhook(url, addr))
	if err != nil {
		t.Errorf("Error creating server with webhook: %s", err)
	}
	if server.webhookURL != url {
		t.Error("Server webhookURL should be set")
	}
	if server.listenAddr != addr {
		t.Error("Server listenAddr should be set")
	}
}

func TestAddMiddleware(t *testing.T) {
	server := &Server{}
	if len(server.middlewares) > 0 {
		t.Error("Middleware list should be empty by default")
	}
	server.AddMiddleware(func(HandlerFunction) HandlerFunction { return nil })
	if len(server.middlewares) != 1 {
		t.Error("AddMiddleware should add new middleware")
	}
}
