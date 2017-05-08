package tbot

import "testing"

func TestDefaultMux(t *testing.T) {
	mux := NewDefaultMux()
	mux.HandleFunc("hi", func(m *Message) { m.Reply("hi") }, "hi")
	handler, _ := mux.Mux("hi")
	if handler == nil {
		t.Fail()
	}
}

func TestReplaceVariables(t *testing.T) {
	pattern := "/say {text}"
	regex := replaceVariables(pattern)
	if regex != "/say (.*)" {
		t.Fail()
	}
}

func TestDefaultMuxWithVariable(t *testing.T) {
	mux := NewDefaultMux()
	mux.HandleFunc("/say {text}", func(m *Message) { m.Reply("hi") })
	handler, data := mux.Mux("/say hi")
	if handler == nil {
		t.Fail()
	}
	if data["text"] != "hi" {
		t.Fail()
	}
}

func TestDefaultMuxWithVariables(t *testing.T) {
	mux := NewDefaultMux()
	mux.HandleFunc("/say {some} {text}", func(m *Message) { m.Reply("hi") })
	_, data := mux.Mux("/say something new")
	if data["some"] != "something" {
		t.Fail()
	}
	if data["text"] != "new" {
		t.Fail()
	}
}

func TestParseVariables(t *testing.T) {
	pattern := "some pattern with {command}"
	vars := parseVariables(pattern)
	if len(vars) != 1 {
		t.FailNow()
	}
	if vars[0] != "command" {
		t.Fail()
	}
}

func TestMuxDefaultHandler(t *testing.T) {
	mux := NewDefaultMux()
	f := func(m *Message) { m.Reply("default") }
	mux.HandleDefault(f)
	handler, err := mux.Mux("some text here")
	if err != nil {
		t.Fail()
	}
	if handler == nil {
		t.Fail()
	}
}

func TestDefaultHandler(t *testing.T) {
	mux := NewDefaultMux()
	mux.HandleDefault(func(m *Message) {})
	handler := mux.DefaultHandler()
	if handler == nil {
		t.Fail()
	}
}

func TestHandlers(t *testing.T) {
	mux := NewDefaultMux()
	mux.HandleFunc("/hi", func(m *Message) {})
	mux.HandleFunc("/test", func(m *Message) {})
	handlers := mux.Handlers()
	if len(handlers) != 2 {
		t.Fail()
	}
	if handlers["/hi"] == nil {
		t.Fail()
	}
}
