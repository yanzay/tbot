package tbot

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/yanzay/tbot/internal/adapter"
	"github.com/yanzay/tbot/model"
)

func TestMain(m *testing.M) {
	createBot = func(token string, httpClient *http.Client) (adapter.BotAdapter, error) {
		if token == TestToken {
			return &mockBot{}, nil
		}
		return nil, fmt.Errorf("wrong token")
	}
	os.Exit(m.Run())
}

var inputMessages = make(chan *model.Message)
var outputMessages = make(chan *model.Message)

type mockBot struct{}

func (*mockBot) Send(m *model.Message) error {
	go func() {
		outputMessages <- m
	}()
	return nil
}

func (*mockBot) SendRaw(endpoint string, params map[string]string) error {
	return nil
}

func (*mockBot) GetUpdatesChan(webhookURL string, listenAddr string) (<-chan *model.Message, error) {
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
	requestResponse(t, setup, "/hi", model.MessageText, "hi", model.MessageText)
}

func TestStickerReply(t *testing.T) {
	setup := func(s *Server) {
		s.HandleFunc("/sticker", func(m *Message) { m.ReplySticker("sticker.png") })
	}
	requestResponse(t, setup, "/sticker", model.MessageText, "sticker.png", model.MessageSticker)
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
		model.MessageDocument, "OK", model.MessageText)
}

func TestKeyboardReply(t *testing.T) {
	setup := func(s *Server) {
		s.HandleFunc("/keyboard", func(m *Message) {
			m.ReplyKeyboard("keys", [][]string{{"hi"}})
		})
	}
	requestResponse(t, setup, "/keyboard", model.MessageText, "keys", model.MessageKeyboard)
}

func requestResponse(t *testing.T, setup func(*Server), inData string, inType model.MessageType, outData string, outType model.MessageType) {
	inputMessages = make(chan *model.Message)
	defer close(inputMessages)
	outputMessages = make(chan *model.Message)
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
	inputMessages <- &model.Message{
		Data:    inData,
		Type:    inType,
		Replies: make(chan *model.Message),
	}
	out := <-outputMessages
	if out.Type != outType {
		t.Errorf("Output should be text")
	}
	if out.Data != outData {
		t.Errorf("Output should be hi")
	}
}
