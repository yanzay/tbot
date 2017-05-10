package tbot

import (
	"fmt"
	"os"
	"testing"

	"github.com/yanzay/tbot/internal/adapter"
)

func TestMain(m *testing.M) {
	createBot = func(token string) (adapter.BotAdapter, error) {
		if token == TestToken {
			return &mockBot{}, nil
		}
		return nil, fmt.Errorf("wrong token")
	}
	os.Exit(m.Run())
}

var inputMessages = make(chan *adapter.Message)
var outputMessages = make(chan *adapter.Message)

type mockBot struct{}

func (*mockBot) Send(m *adapter.Message) error {
	go func() {
		outputMessages <- m
	}()
	return nil
}

func (*mockBot) GetUpdatesChan() (<-chan *adapter.Message, error) {
	return inputMessages, nil
}

func (*mockBot) GetUserName() string {
	return "TestUserName"
}

func (*mockBot) GetFirstName() string {
	return "TestFirstName"
}

func TestTextReply(t *testing.T) {
	setup := func(s *Server) {
		s.Handle("/hi", "hi")
	}
	requestResponse(t, setup, "/hi", adapter.MessageText, "hi", adapter.MessageText)
}

func TestStickerReply(t *testing.T) {
	setup := func(s *Server) {
		s.HandleFunc("/sticker", func(m *Message) { m.ReplySticker("sticker.png") })
	}
	requestResponse(t, setup, "/sticker", adapter.MessageText, "sticker.png", adapter.MessageSticker)
}

func TestDocumentUpload(t *testing.T) {
	setup := func(s *Server) {
		s.HandleFile(func(m *Message) {
			err := m.Download("uploads")
			if err != nil {
				t.Errorf("Error downloading file: %q", err)
			}
			m.Reply("OK")
		})
	}
	requestResponse(t, setup,
		"https://raw.githubusercontent.com/yanzay/tbot/master/LICENSE",
		adapter.MessageDocument, "OK", adapter.MessageText)
}

func requestResponse(t *testing.T, setup func(*Server), inData string, inType adapter.MessageType, outData string, outType adapter.MessageType) {
	inputMessages = make(chan *adapter.Message)
	defer close(inputMessages)
	outputMessages = make(chan *adapter.Message)
	defer close(outputMessages)
	s, err := NewServer("TEST:TOKEN")
	if err != nil {
		t.Errorf("Error creating server: %q", err)
	}
	setup(s)
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			t.Errorf("Error listening and serving: %q", err)
		}
	}()
	inputMessages <- &adapter.Message{
		Data:    inData,
		Type:    inType,
		Replies: make(chan *adapter.Message),
	}
	out := <-outputMessages
	if out.Type != outType {
		t.Errorf("Output should be text")
	}
	if out.Data != outData {
		t.Errorf("Output should be hi")
	}
}
