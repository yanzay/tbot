package tbot

import "testing"

func TestProcessMessage(t *testing.T) {
	server, _ := NewServer("TEST:TOKEN")
	message := mockMessage()
	server.processMessage(message)
	select {
	case <-message.Replies:
		t.Error("should not be replied")
	default:
	}
}
