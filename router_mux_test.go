package tbot

import (
	"strings"
	"testing"

	"github.com/yanzay/tbot/model"
)

type testSequence struct {
	flow     []string
	expected string
}

func TestRouterMux(t *testing.T) {
	seqs := []testSequence{
		{
			[]string{RouteRoot, "pets", "cat"},
			"index pets cat",
		},
		{
			[]string{"pets", RouteBack},
			"pets index",
		},
		{
			[]string{"pets", "cat", RouteBack},
			"pets cat index",
		},
		{
			[]string{"pets", "cat", RouteRoot, "pets"},
			"pets cat index pets",
		},
		{
			[]string{"pets", RouteRefresh},
			"pets pets",
		},
		{
			[]string{"meals", "pizza", "popcorn", RouteBack},
			"meals pizza popcorn index",
		},
		{
			[]string{"/pets", "cat", RouteRoot, "pets"},
			"pets cat index pets",
		},
	}
	for _, seq := range seqs {
		rm := NewRouterMux(NewSessionStorage())
		rm.HandleFunc(RouteRoot, indexHandler)
		rm.HandleFunc("/start", startHandler)
		rm.HandleFunc("/start {var}", startVarHandler)
		rm.HandleFunc("/pets", petsHandler)
		rm.HandleFunc("/pets/cat", catHandler)
		rm.HandleFunc("/pets/cat {var}", catVarHandler)
		rm.HandleFunc("/meals", mealsHandler)
		rm.HandleFunc("/meals/pizza", pizzaHandler)
		rm.HandleFunc("/meals/popcorn", popcornHandler)
		routerMuxFlow(t, rm, seq)
	}
}

func TestRouterAliases(t *testing.T) {
	seqs := []testSequence{
		{
			[]string{"Home", "Pets", "Cat"},
			"index pets cat",
		},
		{
			[]string{"Home", "Pets", RouteBack, "Pictures", "Cat"},
			"index pets index pictures piccat",
		},
		{
			[]string{"Home", "pets", "Kitty", RouteRefresh},
			"index pets cat pets",
		},
		{
			[]string{"/start", "pets", "Kitty meow"},
			"index pets meow",
		},
		{
			[]string{"/start 128325", "pets", "Cat"},
			"128325 pets cat",
		},
		{
			[]string{"  /start 128325", "pets", "Cat"},
			"128325 pets cat",
		}, {
			[]string{"/start 128 325", "pets", "Cat", "unknown-text"},
			"128 325 pets cat default",
		},
	}
	for _, seq := range seqs {
		rm := NewRouterMux(NewSessionStorage())
		rm.HandleDefault(defaultHandler)
		rm.HandleFunc(RouteRoot, indexHandler)
		rm.HandleFunc("/start", startHandler)
		rm.HandleFunc("/start {var}", startVarHandler)
		rm.HandleFunc("/pets", petsHandler)
		rm.HandleFunc("/pets/cat", catHandler)
		rm.HandleFunc("/pets/cat {var}", catVarHandler)
		rm.HandleFunc("/pictures", pictureshandler)
		rm.HandleFunc("/pictures/cat", picCatHandler)
		rm.SetAlias(RouteRoot, "Home")
		rm.SetAlias("pets", "Pets")
		rm.SetAlias("cat", "Cat", "Kitty")
		rm.SetAlias("pictures", "Pictures")
		routerMuxFlow(t, rm, seq)
	}
}

func routerMuxFlow(t *testing.T, mux Mux, seq testSequence) {
	path := make([]string, 0)
	for _, input := range seq.flow {
		msg := &Message{
			Message: &model.Message{Data: input},
			Vars:    make(map[string]string),
		}
		h,vars := mux.Mux(msg)
		if vars != nil {
			msg.Vars = vars
		}
		if h == nil {
			t.Errorf("Handler is nil for message: %s", input)
		}
		h.f(msg)
		path = append(path, msg.Vars["path"])
	}
	fullPath := strings.Join(path, " ")
	if fullPath != seq.expected {
		t.Errorf("Expected path: '%s', actual path: '%s'", seq.expected, fullPath)
	}
}

func defaultHandler(m *Message)  { m.Vars["path"] = "default" }
func indexHandler(m *Message)    { m.Vars["path"] = "index" }
func startHandler(m *Message)    { m.Vars["path"] = "index" }
func startVarHandler(m *Message) { m.Vars["path"] = m.Vars["var"] }
func petsHandler(m *Message)     { m.Vars["path"] = "pets" }
func catHandler(m *Message)      { m.Vars["path"] = "cat" }
func catVarHandler(m *Message)   { m.Vars["path"] = m.Vars["var"] }
func pictureshandler(m *Message) { m.Vars["path"] = "pictures" }
func picCatHandler(m *Message)   { m.Vars["path"] = "piccat" }
func mealsHandler(m *Message)    { m.Vars["path"] = "meals" }
func pizzaHandler(m *Message)    { m.Vars["path"] = "pizza" }
func popcornHandler(m *Message)  { m.Vars["path"] = "popcorn" }
