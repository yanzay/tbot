package tbot_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yanzay/tbot/v2"
)

const token = "TOKEN"

func TestNewClient(t *testing.T) {
	c := tbot.NewClient(token, nil, "https://example.com")
	if c == nil {
		t.Fatalf("client is nil")
	}
}

func TestGetMe(t *testing.T) {
	c := testClient(t, `
		{
			"ok": true,
			"result": {"id": 1}
		}
	`)

	me, err := c.GetMe()
	if err != nil {
		t.Fatalf("error on getMe: %v", err)
	}
	if me.ID == 0 {
		t.Fatalf("empty me.ID")
	}
}

func TestSendMessage(t *testing.T) {
	c := testClient(t, `
		{
			"result": {
				"chat": {"id": 1},
				"text": "helo"
			},
			"ok": true
		}
	`)

	msg, err := c.SendMessage("123", "helo")
	if err != nil {
		t.Fatalf("error on sendMessage: %v", err)
	}
	if msg.Text == "" {
		t.Fatalf("empty message text")
	}
}

func TestSendMessageWithOptions(t *testing.T) {
	c := testClient(t, `
		{
			"result": {
				"chat": {"id": 1},
				"text": "helo"
			},
			"ok": true
		}
	`)

	msg, err := c.SendMessage("123", "helo", tbot.OptParseModeMarkdown,
		tbot.OptDisableWebPagePreview, tbot.OptDisableNotification,
		tbot.OptReplyToMessageID(1), tbot.OptForceReply, tbot.OptReplyKeyboardRemove)
	if err != nil {
		t.Fatalf("error on sendMessage: %v", err)
	}
	if msg.Text == "" {
		t.Fatalf("empty message text")
	}
}

func TestForwardMessage(t *testing.T) {
	c := testClient(t, `
		{
			"ok": true,
			"result": {"message_id": 321}
		}
	`)

	msg, err := c.ForwardMessage("321", "123", 1)
	if err != nil {
		t.Fatalf("error on forwardMessage: %v", err)
	}
	if msg.MessageID == 0 {
		t.Fatalf("empty message id")
	}
}

func TestSendAudio(t *testing.T) {
	c := testClient(t, `
		{
			"ok": true,
			"result": {"message_id": 321}
		}
	`)
	msg, err := c.SendAudio("123", "aaa")
	if err != nil {
		t.Fatalf("error on sendAudio: %v", err)
	}
	if msg.MessageID == 0 {
		t.Fatalf("empty message id")
	}
}

func TestSendAudioFile(t *testing.T) {
	c := testClient(t, `
		{
			"ok": true,
			"result": {"message_id": 321}
		}
	`)
	msg, err := c.SendAudioFile("123", "client_test.go")
	if err != nil {
		t.Fatalf("error on sendAudioFile: %v", err)
	}
	if msg.MessageID == 0 {
		t.Fatalf("empty message id")
	}
}

func TestSendPhoto(t *testing.T) {
	c := testClient(t, `
		{
			"ok": true,
			"result": {"message_id": 321}
		}
	`)
	msg, err := c.SendPhoto("123", "aaa")
	if err != nil {
		t.Fatalf("error on sendPhoto: %v", err)
	}
	if msg.MessageID == 0 {
		t.Fatalf("empty message id")
	}
}

func TestSendPhotoFile(t *testing.T) {
	c := testClient(t, `
		{
			"ok": true,
			"result": {"message_id": 321}
		}
	`)
	msg, err := c.SendPhotoFile("123", "client_test.go")
	if err != nil {
		t.Fatalf("error on sendPhotoFile: %v", err)
	}
	if msg.MessageID == 0 {
		t.Fatalf("empty message id")
	}
}

func TestSendDice(t *testing.T) {
	c := testClient(t, `
		{
			"ok": true,
			"result": {
				"emoji": "🎲",
				"value": 6
			}
		}
	`)
	msg, err := c.SendDice("123", "🎲")
	if err != nil {
		t.Fatalf("error on sendDice: %v", err)
	}
	if msg.Value == 0 {
		t.Fatalf("empty dice value")
	}
}

func testClient(t *testing.T, resp string) *tbot.Client {
	t.Helper()
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, resp)
	}
	httpServer := httptest.NewServer(http.HandlerFunc(handler))
	httpClient := httpServer.Client()
	return tbot.NewClient(token, httpClient, httpServer.URL)
}
