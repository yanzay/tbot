package tbot

import (
	"fmt"
	"os"
	"testing"

	"github.com/yanzay/tbot/adapter"
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

type mockBot struct{}

func (*mockBot) Send(*adapter.Message) error {
	return nil
}

func (*mockBot) GetUpdatesChan() (<-chan *adapter.Message, error) {
	return make(chan *adapter.Message), nil
}

func (*mockBot) GetFileDirectURL(string) (string, error) {
	return "", nil
}

func (*mockBot) GetUserName() string {
	return "TestUserName"
}

func (*mockBot) GetFirstName() string {
	return "TestFirstName"
}
