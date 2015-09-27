package tbot

import "testing"

func TestDefaultMux(t *testing.T) {
	mux := NewDefaultMux()
	mux.HandleFunc("hi", func(m Message) { m.Reply("hi") }, "hi")
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
	mux.HandleFunc("/say {text}", func(m Message) { m.Reply("hi") })
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
	mux.HandleFunc("/say {some} {text}", func(m Message) { m.Reply("hi") })
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
