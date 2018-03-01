package tbot

import (
	"testing"
)

func TestAccessDenied(t *testing.T) {
	message := mockMessage()
	go accessDenied(message, "Userid")
	reply := <-message.Replies
	if reply == nil {
		t.Fail()
	}
}

func TestAuthWithUserNameSuccess(t *testing.T) {
	auth := NewAuth([]string{"me"})
	message := mockMessage()
	invoked := false
	handler := func(m *Message) { invoked = true }
	auth(handler)(message)
	if !invoked {
		t.Fail()
	}
}

func TestAuthWithUserIdSuccess(t *testing.T) {
	auth := NewAuth([]int{521844285})
	message := mockMessage()
	invoked := false
	handler := func(m *Message) { invoked = true }
	auth(handler)(message)
	if !invoked {
		t.Fail()
	}
}

func TestAuthWithChatIdSuccess(t *testing.T) {
	auth := NewAuth([]int64{-271651159})
	message := mockMessage()
	invoked := false
	handler := func(m *Message) { invoked = true }
	auth(handler)(message)
	if !invoked {
		t.Fail()
	}
}

func TestAuthUserNameFail(t *testing.T) {
	auth := NewAuth([]string{"notme"})
	message := mockMessage()
	invoked := false
	handler := func(m *Message) { invoked = true }
	go auth(handler)(message)
	<-message.Replies
	if invoked {
		t.Fail()
	}
}

func TestAuthUserIdFail(t *testing.T) {
	auth := NewAuth([]int{1234567})
	message := mockMessage()
	invoked := false
	handler := func(m *Message) { invoked = true }
	go auth(handler)(message)
	<-message.Replies
	if invoked {
		t.Fail()
	}
}

func TestAuthChatIdFail(t *testing.T) {
	auth := NewAuth([]int64{11235813})
	message := mockMessage()
	invoked := false
	handler := func(m *Message) { invoked = true }
	go auth(handler)(message)
	<-message.Replies
	if invoked {
		t.Fail()
	}
}
