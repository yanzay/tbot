package tbot

import (
	"errors"
	"fmt"
)

const (
	username = "Username"
	userid   = "Userid"
	chatid   = "Chatid"
)

//NewAuth creates Middleware for white-list based
//authentication according to username, userid or chatid list.
//purpose: to prevent the access to bots from another users or groups
func NewAuth(whitelist interface{}) Middleware {
	switch whitelist.(type) {
	case []string:
		return NewAuthWithUserName(whitelist.([]string))
	case []int:
		return NewAuthWithUserId(whitelist.([]int))
	case []int64:
		return NewAuthWithChatId(whitelist.([]int64))
	default:
		panic(errors.New("Unknown Whitelist Format"))
	}
}

func NewAuthWithUserName(whitelist []string) Middleware {
	return func(f HandlerFunction) HandlerFunction {
		return func(m *Message) {
			for _, name := range whitelist {
				fmt.Println(m.From.ID)
				if m.From.UserName == name {
					f(m)
					return
				}
			}
			accessDenied(m, username)
		}
	}
}

func NewAuthWithUserId(whitelist []int) Middleware {
	return func(f HandlerFunction) HandlerFunction {
		return func(m *Message) {
			for _, userId := range whitelist {
				if m.From.ID == userId {
					f(m)
					return
				}
			}
			accessDenied(m, userid)
		}
	}
}

func NewAuthWithChatId(whitelist []int64) Middleware {
	return func(f HandlerFunction) HandlerFunction {
		return func(m *Message) {
			for _, chatId := range whitelist {
				if m.ChatID == chatId {
					f(m)
					return
				}
			}
			accessDenied(m, chatid)
		}
	}
}

func accessDenied(m *Message, opt string) {
	switch opt {
	case username:
		m.Replyf("Access Denied For This %s: %s", opt, m.From.UserName)
	case userid:
		m.Replyf("Access Denied For This %s: %d", opt, m.From.ID)
	case chatid:
		m.Replyf("Access Denied For This %s: %d", opt, m.ChatID)
	default:
	}
}
