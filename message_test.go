package tbot

import (
	"testing"

	"github.com/yanzay/tbot/adapter"
)

func TestReply(t *testing.T) {
	message := mockMessage()
	go message.Reply("hi there")
	reply := <-message.Replies
	if reply.Data != "hi there" {
		t.Fail()
	}
}

func TestReplyf(t *testing.T) {
	message := mockMessage()
	go message.Replyf("the answer is %d", 42)
	reply := <-message.Replies
	if reply.Data != "the answer is 42" {
		t.Fail()
	}
}

func TestReplySticker(t *testing.T) {
	message := mockMessage()
	go message.ReplySticker("server.go")
	reply := <-message.Replies
	if reply.Data == "" {
		t.Error("Reply should contain sticker url")
	}
}

func TestReplyPhoto(t *testing.T) {
	message := mockMessage()
	go message.ReplyPhoto("server.go", "it's me")
	reply := <-message.Replies
	if reply.Data == "" {
		t.Error("Reply should contain photo url")
	}
}

func TestReplyAudio(t *testing.T) {
	message := mockMessage()
	go message.ReplyAudio("server.go")
	reply := <-message.Replies
	if reply.Data == "" {
		t.Error("Reply should contain audio url")
	}
}

func TestReplyDocument(t *testing.T) {
	message := mockMessage()
	go message.ReplyDocument("server.go")
	reply := <-message.Replies
	if reply.Data == "" {
		t.Error("Reply should contain document url")
	}
}

func mockMessage() *Message {
	m := &Message{
		Message: &adapter.Message{
			ChatID:  13666,
			From:    "me",
			Replies: make(chan *adapter.Message),
		},
	}
	return m
}
