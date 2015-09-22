package tbot

import "testing"

func TestDefaultMux(t *testing.T) {
	handler := NewHandler(func(m Message) { m.Reply("hi") }, "hi")
	handlers := map[string]*Handler{
		"hi": handler,
	}
	handler, _ = DefaultMux(handlers, "hi")
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
	handler := NewHandler(func(m Message) { m.Reply("hi") }, "/say {text}")
	handlers := map[string]*Handler{
		"/say {text}": handler,
	}
	handler, data := DefaultMux(handlers, "/say hi")
	if handler == nil {
		t.Fail()
	}
	if data["text"] != "hi" {
		t.Log("data[text]: " + data["text"])
		t.Fail()
	}
}

func TestDefaultMuxWithVariables(t *testing.T) {
	handler := NewHandler(func(m Message) { m.Reply("hi") }, "/say {some} {text}")
	handlers := map[string]*Handler{
		"/say {some} {text}": handler,
	}
	_, data := DefaultMux(handlers, "/say something new")
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
