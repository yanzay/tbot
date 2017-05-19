package tbot

import (
	"strings"
	"testing"

	"github.com/yanzay/tbot/internal/adapter"
)

type testSequence struct {
	flow     []string
	expected string
}

func TestRouterMux(t *testing.T) {
	seqs := []testSequence{
		{
			[]string{"/index", "/pets", "/cat"},
			"index pets cat",
		},
		{
			[]string{"/index", "/pets", RouteBack},
			"index pets index",
		},
		{
			[]string{"/index", "/pets", "/cat", RouteBack, RouteBack},
			"index pets cat pets index",
		},
		{
			[]string{"/index", "/pets", "/cat", RouteRoot, "/pets"},
			"index pets cat index pets",
		},
		{
			[]string{"/index", "/pets", RouteRefresh},
			"index pets pets",
		},
	}
	for _, seq := range seqs {
		routerMuxFlow(t, seq.flow, seq.expected)
	}
}

func routerMuxFlow(t *testing.T, flow []string, expected string) {
	sessions := NewSessionStorage()
	rm := NewRouterMux(sessions)
	path := make([]string, 0)
	indexHandler := func(*Message) { path = append(path, "index") }
	petsHandler := func(*Message) { path = append(path, "pets") }
	catHandler := func(*Message) { path = append(path, "cat") }
	rm.HandleFunc("/index", indexHandler)
	rm.HandleFunc("/index/pets", petsHandler)
	rm.HandleFunc("/index/pets/cat", catHandler)
	for _, input := range flow {
		msg := &Message{Message: &adapter.Message{Data: input}}
		h, _ := rm.Mux(msg)
		if h == nil {
			t.Errorf("Handler is nil for message: %s", input)
		}
		h.f(msg)
	}
	fullPath := strings.Join(path, " ")
	if fullPath != expected {
		t.Errorf("Expected path: '%s', actual path: '%s'", expected, fullPath)
	}
}
