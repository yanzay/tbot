package tbot

import (
	"fmt"
	"os"
	"testing"

	"github.com/yanzay/tbot/adapter"
)

func TestMain(m *testing.M) {
	createBot = func(token string) (*adapter.Bot, error) {
		if token == TestToken {
			return &adapter.Bot{}, nil
		}
		return nil, fmt.Errorf("Wrong token!")
	}
	os.Exit(m.Run())
}
