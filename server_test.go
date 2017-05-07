package tbot

import "testing"

const (
	TestToken    = "153667468:AAHlSHlMqSt1f_uFmVRJbm5gntu2HI4WW8I"
	InvalidToken = "invalid"
)

func TestNewServerSuccess(t *testing.T) {
	server, err := NewServer(TestToken)
	if err != nil {
		t.Fail()
	}
	if server == nil {
		t.Fail()
	}
	if server.mux == nil {
		t.Fail()
	}
}

func TestNewServerFail(t *testing.T) {
	server, err := NewServer(InvalidToken)
	if err == nil {
		t.Fail()
	}
	if server != nil {
		t.Fail()
	}
}

func TestAddMiddleware(t *testing.T) {
	server := &Server{}
	if len(server.middlewares) > 0 {
		t.Fail()
	}
	server.AddMiddleware(func(HandlerFunction) HandlerFunction { return nil })
	if len(server.middlewares) != 1 {
		t.Fail()
	}
}
