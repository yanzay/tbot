package tbot

import "testing"

func TestNewHandler(t *testing.T) {
	handler := NewHandler(func(*Message) {}, "/", "Desc")
	if handler == nil {
		t.Error("NewHandler should create non-nil handler")
	}
}
