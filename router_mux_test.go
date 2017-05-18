package tbot

import (
	"strings"
	"testing"

	"github.com/yanzay/tbot/internal/adapter"
)

func TestRouterMux(t *testing.T) {
	flows := [][]string{
		{"/index", "/pets", "/cat"},
		{"/index", "/pets", RouteBack},
		{"/index", "/pets", "/cat", RouteBack, RouteBack},
		{"/index", "/pets", "/cat", RouteRoot, "/pets"},
	}
	expected := []string{
		"index pets cat",
		"index pets index",
		"index pets cat pets index",
		"index pets cat index pets",
	}
	for i := range flows {
		routerMuxFlow(t, flows[i], expected[i])
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
