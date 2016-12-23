package tbot

import "testing"

func TestHelpHandler(t *testing.T) {
	s := &Server{
		mux: NewDefaultMux(),
	}
	s.HandleFunc("/test", func(m Message) {}, "desc")
	s.HandleDefault(func(m Message) {}, "default desc")
	message := mockMessage()
	go func() { s.HelpHandler(message) }()
	reply := <-message.replies
	if reply.msg == nil {
		t.Fail()
	}
}
