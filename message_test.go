package tbot

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func TestReply(t *testing.T) {
	message := mockMessage()
	go message.Reply("hi there")
	reply := <-message.replies
	replyMsg, ok := reply.msg.(tgbotapi.MessageConfig)
	if !ok {
		t.Fail()
	}
	if replyMsg.Text != "hi there" {
		t.Fail()
	}
}

func TestReplyf(t *testing.T) {
	message := mockMessage()
	go message.Replyf("the answer is %d", 42)
	reply := <-message.replies
	if reply.msg.(tgbotapi.MessageConfig).Text != "the answer is 42" {
		t.Fail()
	}
}

func TestReplySticker(t *testing.T) {
	message := mockMessage()
	go message.ReplySticker("server.go")
	reply := <-message.replies
	_, ok := reply.msg.(tgbotapi.StickerConfig)
	if !ok {
		t.Fail()
	}
}

func TestReplyPhoto(t *testing.T) {
	message := mockMessage()
	go message.ReplyPhoto("server.go", "it's me")
	reply := <-message.replies
	_, ok := reply.msg.(tgbotapi.PhotoConfig)
	if !ok {
		t.Fail()
	}
}

func TestReplyAudio(t *testing.T) {
	message := mockMessage()
	go message.ReplyAudio("server.go")
	reply := <-message.replies
	_, ok := reply.msg.(tgbotapi.AudioConfig)
	if !ok {
		t.Fail()
	}
}

func TestReplyDocument(t *testing.T) {
	message := mockMessage()
	go message.ReplyDocument("server.go")
	reply := <-message.replies
	_, ok := reply.msg.(tgbotapi.DocumentConfig)
	if !ok {
		t.Fail()
	}
}

func mockMessage() Message {
	chat := &tgbotapi.Chat{}
	user := &tgbotapi.User{UserName: "me"}
	m := Message{tgbotapi.Message{Chat: chat, From: user}, map[string]string{}, make(chan *ReplyMessage)}
	return m
}
